package store

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/im-kulikov/simplinic-task/models"
	"github.com/pkg/errors"
)

func (s *schemes) Create(scheme *models.Scheme) error {
	//scheme.ID = 0 // set to null

	if _, err := s.db.Model(scheme).
		Insert(); err != nil {
		return errors.WithMessage(err, "can't create scheme")
	}

	return nil
}

func (s *schemes) Read(id int64) (*models.Scheme, error) {
	var result models.Scheme

	if err := s.db.Model(&result).
		Where("id = ?", id).First(); err != nil {
		return nil, errors.Wrapf(err, "can't read scheme #%d", id)
	}

	return &result, nil
}

func (s *schemes) Update(scheme *models.Scheme) error {
	var version int64

	if err := s.db.Model((*models.Scheme)(nil)).
		Column("version").
		Where("id = ?", scheme.ID).
		Limit(1).
		Select(pg.Scan(&version)); err != nil {
		return errors.Wrapf(err, "can't fetch version for scheme #%d", scheme.ID)
	}

	scheme.ID = 0 // drop id
	scheme.Version = version + 1

	_, err := s.db.Model(scheme).
		Insert()

	return errors.WithMessage(err, "can't create scheme")
}

func (s *schemes) Delete(scheme *models.Scheme) error {
	scheme.DeletedAt.Time = time.Now()

	if _, err := s.db.Model(scheme).
		Column("deleted_at").
		Where("id = ?", scheme.ID).
		Update(); err != nil {
		return errors.Wrapf(err, "can't remove scheme #%d", scheme.ID)
	}

	return nil
}

func (s *schemes) Search(req *SearchRequest) ([]*models.Scheme, error) {
	var result []*models.Scheme

	q := s.db.Model(&result)

	if req.Version > 0 {
		q.Where("version = ?", req.Version)
	}

	if len(req.Tags) > 0 {
		q.Where(`tags @> ?`, req.Tags) // tags @> '["b", "c"]' : filter tags, that have "b" and "c"
	}

	q.Where("deleted_at ISNULL")

	if err := q.Select(); err != nil {
		return nil, errors.Wrapf(err, "can't find schemes by (version=%d | tags=%v)", req.Version, req.Tags)
	}

	return result, nil
}
