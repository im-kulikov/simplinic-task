package store

import (
	"encoding/json"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-pg/pg"
	"github.com/im-kulikov/simplinic-task/models"
	"github.com/pkg/errors"
)

type (
	Scheme struct {
		tableName struct{}        `sql:"scheme_versions,alias:sv" pg:",discard_unknown_columns"`
		ID        int64           `sql:"scheme_id" json:"id"`
		Version   int64           `json:"version"`
		Tags      []string        `json:"tags"`
		Data      json.RawMessage `json:"data"`
	}
)

func (s *schemes) Create(scheme *Scheme) error {
	var model models.Scheme

	if _, err := s.db.Model(&model).Insert(); err != nil {
		return errors.WithMessage(err, "could not create scheme")
	}

	scheme.ID = model.ID

	if _, err := s.db.Model(scheme).Insert(); err != nil {
		spew.Dump(err)
		return errors.WithMessage(err, "could not create scheme data")
	}

	return nil
}

func (s *schemes) Read(id int64) (*Scheme, error) {
	var result Scheme

	if err := s.db.Model(&result).
		Where("scheme_id = ?", id).
		Order("created_at DESC", "version DESC").
		Limit(1).
		Select(); err != nil {
		return nil, errors.Wrapf(err, "could not read scheme #%d", id)
	}

	return &result, nil
}

func (s *schemes) Update(scheme *Scheme) error {
	var id, version int64

	if err := s.db.Model((*models.Scheme)(nil)).
		ColumnExpr("id, version").
		Join("LEFT JOIN scheme_versions sv").
		JoinOn("sv.scheme_id = scheme.id").
		Where("id = ? AND deleted_at ISNULL", scheme.ID).
		Order("sv.version DESC").
		Limit(1).
		Select(pg.Scan(&id, &version)); err != nil {
		return errors.Wrapf(err, "could not update scheme #%d", scheme.ID)
	} else if id != scheme.ID {
		return errors.Errorf("could not update scheme #%d, not found", scheme.ID)
	}

	scheme.Version = version + 1

	if _, err := s.db.Model(scheme).
		Insert(); err != nil {
		return errors.WithMessage(err, "can't create scheme")
	}

	//return
	return nil
}

func (s *schemes) Delete(id int64) error {
	if err := s.db.Delete(&models.Scheme{ID: id}); err != nil {
		return errors.Wrapf(err, "can't remove scheme #%d", id)
	}

	return nil
}

func (s *schemes) Search(req SearchRequest) ([]*Scheme, error) {
	var result []*Scheme

	q := s.db.Model(&result).
		ColumnExpr("sv.*").
		Join("LEFT JOIN schemes s").
		JoinOn("s.id = sv.scheme_id").
		Order("version DESC").
		Group("scheme_id", "version")

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
