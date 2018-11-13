package models

import (
	"encoding/json"

	"github.com/go-pg/pg"
)

type Config struct {
	ID        int64 `pg:",pk"`
	Version   int64
	SchemesID int64
	Tags      []string
	Data      json.RawMessage
	CreatedAt pg.NullTime
	DeletedAt pg.NullTime
}
