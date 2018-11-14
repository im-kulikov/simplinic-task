package store

import (
	"encoding/json"
	"time"

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

func (s *configs) Create(cfg *models.Config) error {
	if _, err := s.db.Model(cfg).
		Insert(); err != nil {
		return errors.WithMessage(err, "can't create config")
	}

	return nil
}

func (s *configs) Read(id int64) (*models.Config, error) {
	var result models.Config

	if err := s.db.Model(&result).
		Where("id = ?", id).First(); err != nil {
		return nil, errors.Wrapf(err, "can't read config #%d", id)
	}

	return &result, nil
}

// update/configs/8 version=1
// update/configs/9 version=2

func (s *configs) Update(cfg *models.Config) error {
	//var version int64
	//
	//if err := s.db.Model((*models.Config)(nil)).
	//	Column("version").
	//	Where("id = ?", cfg.ID).
	//	Limit(1).
	//	Select(pg.Scan(&version)); err != nil {
	//	return errors.Wrapf(err, "can't fetch version for config #%d", cfg.ID)
	//}
	//
	//cfg.ID = 0 // drop id
	//cfg.Version = version + 1
	//
	//_, err := s.db.Model(cfg).
	//	Insert()
	//
	//return errors.WithMessage(err, "can't create config")
	return nil
}

func (s *configs) Delete(cfg *models.Config) error {
	cfg.DeletedAt = time.Now()

	if _, err := s.db.Model(cfg).
		Column("deleted_at").
		Where("id = ?", cfg.ID).
		Update(); err != nil {
		return errors.Wrapf(err, "can't remove config #%d", cfg.ID)
	}

	return nil
}

func (s *configs) Search(req *SearchRequest) ([]*models.Config, error) {
	var result []*models.Config

	q := s.db.Model(&result)

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
