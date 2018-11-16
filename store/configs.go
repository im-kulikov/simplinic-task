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
	SchemeID  int64           `json:"scheme_id" validate:"required,gt=0" message:"scheme_id could not be empty"`
	Version   int64           `json:"version"`
	Tags      []string        `json:"tags" validate:"required" message:"tags could not be empty"`
	Data      json.RawMessage `json:"data" validate:"required" message:"data could not be empty"`
}

func (s *configs) Create(cfg *Config) error {
	var model = models.Config{SchemeID: cfg.ID}

	if _, err := s.db.Model(&model).Insert(); err != nil {
		return errors.WithMessage(err, "could not create config")
	}

	cfg.ID = model.ID

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
	var sid, version int64

	// Example:
	//   SELECT MAX(cv.version) as version
	//     FROM config_versions cv
	//LEFT JOIN configs c
	//       ON cv.config_id = c.id
	//LEFT JOIN schemes s
	//       ON cv.scheme_id = c.scheme_id
	//    WHERE cv.config_id = 3
	//      AND cv.scheme_id = 10
	//      AND c.deleted_at ISNULL
	//      AND s.deleted_at ISNULL
	// GROUP BY c.id;

	if err := s.db.
		Model((*Config)(nil)).
		ColumnExpr("cv.scheme_id, MAX(cv.version) as version").
		Join("LEFT JOIN configs c").JoinOn("c.id = cv.config_id").
		Join("LEFT JOIN schemes s").JoinOn("s.id = cv.scheme_id").
		Where("cv.config_id = ? AND c.deleted_at ISNULL AND s.deleted_at ISNULL", cfg.ID).
		Group("cv.config_id", "cv.scheme_id").
		Select(pg.Scan(&sid, &version)); err != nil {
		return errors.Wrapf(err, "query error for config #%d", cfg.ID)
	} else if version == 0 || sid == 0 {
		return errors.Errorf("could not update config #%d not found", cfg.ID)
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
