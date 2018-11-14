package models

import (
	"encoding/json"
	"time"
)

type Config struct {
	ID        int64     `pg:",pk"`
	SchemeID  int64     `sql:"scheme_id"`
	CreatedAt time.Time `sql:"created_at"`
	DeletedAt time.Time `pg:",soft_delete"`
}

type ConfigVersion struct {
	ConfigID  int64
	Version   int64
	Tags      []string
	Data      json.RawMessage
	CreatedAt time.Time
}
