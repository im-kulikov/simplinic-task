package store

import (
	"encoding/json"

	"github.com/go-pg/pg"
	"github.com/im-kulikov/simplinic-task/models"
	"github.com/pkg/errors"
)

type Config struct {
	tableName struct{}        `sql:"config_versions,alias:cv" pg:",discard_unknown_columns"`
	ID        int64           `sql:"config_id" json:"id"`
	SchemeID  int64           `json:"scheme_id"`
	Version   int64           `json:"version"`
	Tags      []string        `json:"tags"`
	Data      json.RawMessage `json:"data"`
}

func (s *configs) Create(cfg *Config) error {
	var (
		id, version int64
		model       = models.Config{SchemeID: cfg.SchemeID}
	)

	if err := s.db.
		Model((*models.Scheme)(nil)).
		Column("c.id", "cv.version").
		Join("LEFT JOIN scheme_versions sv"). // LEFT JOIN scheme_versions sv ON sv.scheme_id = scheme.id
		JoinOn("sv.scheme_id = scheme.id").
		Join("LEFT JOIN configs c"). // LEFT JOIN configs c ON c.scheme_id = scheme.id
		JoinOn("c.id = scheme.id").
		Join("LEFT JOIN config_versions cv"). // LEFT JOIN config_versions cv ON cv.scheme_id = scheme.id
		JoinOn("cv.scheme_id = scheme.id").
		Where("scheme.id = ? AND scheme.deleted_at ISNULL AND c.deleted_at ISNULL", cfg.SchemeID).
		Order("cv.version DESC").
		Limit(1).
		Select(pg.Scan(&id, &version)); err != nil {
		return errors.Wrapf(err, "query error for find scheme #%d", cfg.SchemeID)
	} else if id == 0 { // we not find any config for current scheme
		// that's why we create new record...
		if _, err := s.db.Model(&model).Insert(); err != nil {
			return errors.Wrapf(err, "could not create config (scheme#%d)", cfg.SchemeID)
		}

		// set new config.id
		id = model.ID
	}

	cfg.ID = id               // config_id = id
	cfg.Version = version + 1 // version++

	// create new config_versions..
	if _, err := s.db.Model(cfg).Insert(); err != nil {
		return errors.WithMessage(err, "could not store config_version data")
	}

	return nil
}

func (s *configs) Read(id int64) (*Config, error) {
	var result Config

	if err := s.db.Model(&result).
		Join("LEFT JOIN configs c"). // LEFT JOIN configs c ON c.id = cv.config_id
		JoinOn("c.id = cv.config_id").
		Where("c.id = ? AND c.deleted_at ISNULL", id).
		Order("cv.created_at DESC", "cv.version DESC").
		Limit(1).
		Select(); err != nil {
		return nil, errors.Wrapf(err, "could not read config #%d", id)
	}

	return &result, nil
}

func (s *configs) Update(cfg *Config) error {
	var cid, sid, version int64

	if err := s.db.Model((*models.Config)(nil)).
		Column("config.id", "s.id", "cv.version").
		Join("LEFT JOIN schemes s"). // LEFT JOIN schemes s ON s.id = config.scheme_id
		JoinOn("s.id = config.scheme_id").
		Join("LEFT JOIN config_versions cv"). // LEFT JOIN schemes s ON s.id = config.scheme_id
		JoinOn("cv.scheme_id = config.scheme_id").
		Where("config.id = ? AND config.deleted_at ISNULL AND s.deleted_at ISNULL", cfg.ID).
		Order("cv.version DESC").
		Limit(1).
		Select(pg.Scan(&cid, &sid, &version)); err != nil {
		return errors.Wrapf(err, "query error for config #%d", cfg.ID)
	} else if cid != cfg.ID {
		return errors.Errorf("could not update scheme #%d or config #%d not found", cfg.SchemeID, cfg.ID)
	}

	cfg.SchemeID = sid
	cfg.Version = version + 1

	if _, err := s.db.Model(cfg).
		Insert(); err != nil {
		return errors.WithMessage(err, "could not store new version of config data")
	}

	return nil
}

func (s *configs) Delete(id int64) error {
	if err := s.db.Delete(&models.Config{ID: id}); err != nil {
		return errors.Wrapf(err, "can't remove scheme #%d", id)
	}

	return nil
}

func (s *configs) Search(req SearchRequest) ([]*Config, error) {
	var result []*Config

	q := s.db.Model(&result).
		ColumnExpr("cv.*").
		Join("LEFT JOIN configs c").
		JoinOn("c.id = cv.config_id").
		Order("version DESC").
		Group("cv.scheme_id", "cv.config_id", "cv.version")

	if req.Version > 0 {
		q.Where("version = ?", req.Version)
	}

	if len(req.Tags) > 0 {
		q.Where(`tags @> ?`, req.Tags) // tags @> '["b", "c"]' : filter tags, that have "b" and "c"
	}

	q.Where("deleted_at ISNULL")

	if err := q.Select(); err != nil {
		return nil, errors.Wrapf(err, "can't find config by (version=%d | tags=%v)", req.Version, req.Tags)
	}

	return result, nil
}
