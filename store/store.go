package store

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/im-kulikov/helium/module"
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
		Create(scheme *Scheme) error
		Read(id int64) (*Scheme, error)
		Update(scheme *Scheme) error
		Delete(id int64) error
		Search(req SearchRequest) ([]*Scheme, error)
	}

	Configs interface {
		Create(cfg *Config) error
		Read(id int64) (*Config, error)
		Update(cfg *Config) error
		Delete(id int64) error
		Search(req SearchRequest) ([]*Config, error)
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
