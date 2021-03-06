package store

import (
	"encoding/json"

	"github.com/go-pg/pg"
	"github.com/im-kulikov/simplinic-task/models"
	"github.com/pkg/errors"
)

type (
	Scheme struct {
		tableName struct{}        `sql:"scheme_versions,alias:sv" pg:",discard_unknown_columns"`
		ID        int64           `sql:"scheme_id" json:"id"`
		Version   int64           `json:"version"`
		Tags      []string        `json:"tags" validate:"required" message:"tags could not be empty"`
		Data      json.RawMessage `json:"data" validate:"required" message:"data could not be empty"`
	}
)

func (s *schemes) Create(scheme *Scheme) error {
	var model models.Scheme

	if _, err := s.db.Model(&model).Insert(); err != nil {
		return errors.WithMessage(err, "could not create scheme")
	}

	scheme.ID = model.ID

	if _, err := s.db.Model(scheme).Insert(); err != nil {
		return errors.WithMessage(err, "could not create scheme data")
	}

	return nil
}

func (s *schemes) Read(id int64) (*Scheme, error) {
	var result Scheme

	if err := s.db.Model(&result).
		Join("LEFT JOIN schemes s"). // LEFT JOIN configs c ON c.id = cv.config_id
		JoinOn("s.id = sv.scheme_id").
		Where("s.id = ? AND s.deleted_at ISNULL", id).
		Order("sv.created_at DESC", "sv.version DESC").
		Limit(1).
		Select(); err != nil {
		return nil, errors.Wrapf(err, "could not read scheme #%d", id)
	}

	return &result, nil
}

func (s *schemes) Update(scheme *Scheme) error {
	var id, version int64

	// Example:
	//    SELECT sv.scheme_id, MAX(sv.version)
	//      FROM scheme_versions AS sv
	// LEFT JOIN schemes s
	//        ON s.id = sv.scheme_id
	//     WHERE sv.scheme_id = 3
	//       AND s.deleted_at ISNULL
	//  GROUP BY "sv"."scheme_id"

	if err := s.db.Model((*Scheme)(nil)).
		ColumnExpr("sv.scheme_id, MAX(sv.version)").
		Join("LEFT JOIN schemes s").
		JoinOn("s.id = sv.scheme_id").
		Where("sv.scheme_id = ? AND s.deleted_at ISNULL", scheme.ID).
		Group("sv.scheme_id").
		Select(pg.Scan(&id, &version)); err != nil {
		return errors.Wrapf(err, "query error for scheme #%d", scheme.ID)
	} else if id == 0 {
		return errors.Errorf("could not update scheme #%d, not found", scheme.ID)
	}

	scheme.Version = version + 1

	if _, err := s.db.Model(scheme).
		Insert(); err != nil {
		return errors.WithMessage(err, "can't create scheme")
	}

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
