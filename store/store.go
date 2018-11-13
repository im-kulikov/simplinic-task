package store

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/im-kulikov/helium/module"
	"github.com/im-kulikov/simplinic-task/models"
)

type (
	Store interface {
		Schemes() Schemes
		Configs() Configs
	}

	SearchRequest struct {
		Version int64    `json:"version"`
		Tags    []string `json:"tags"`
	}

	Schemes interface {
		Create(scheme *models.Scheme) error
		Read(id int64) (*models.Scheme, error)
		Update(scheme *models.Scheme) error
		Delete(scheme *models.Scheme) error
		Search(req *SearchRequest) ([]*models.Scheme, error)
	}

	Configs interface {
		Create(cfg *models.Config) error
		Read(id int64) (*models.Config, error)
		Update(cfg *models.Config) error
		Delete(cfg *models.Config) error
		Search(req *SearchRequest) ([]*models.Config, error)
	}

	store struct {
		db *pg.DB
	}

	schemes struct {
		db orm.DB
	}

	configs struct {
		db orm.DB
	}
)

var Module = module.Module{
	{Constructor: newStore},
}

func newStore(db *pg.DB) Store {
	return &store{db: db}
}

func (s *store) Schemes() Schemes {
	return &schemes{db: s.db}
}

func (s *store) Configs() Configs {
	return &configs{db: s.db}
}
