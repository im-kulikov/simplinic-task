package store

import (
	"github.com/go-pg/pg"
	"github.com/im-kulikov/helium/module"
)

type (
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

	schemes struct {
		db *pg.DB
	}

	configs struct {
		db *pg.DB
	}
)

var Module = module.Module{
	{Constructor: NewSchemeStore},
	{Constructor: NewConfigStore},
}

func NewSchemeStore(db *pg.DB) Schemes {
	return &schemes{db: db}
}

func NewConfigStore(db *pg.DB) Configs {
	return &configs{db: db}
}
