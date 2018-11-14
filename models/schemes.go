package models

import (
	"encoding/json"
	"time"
)

type (
	Scheme struct {
		ID        int64     `pg:",pk"`
		CreatedAt time.Time `sql:"created_at"`
		DeletedAt time.Time `sql:"deleted_at" pg:",soft_delete"`
	}

	SchemeVersion struct {
		SchemeID  int64           `json:"scheme_id,omitempty"`
		Version   int64           `json:"version,omitempty"`
		Tags      []string        `json:"tags,omitempty"`
		Data      json.RawMessage `json:"data,omitempty"`
		CreatedAt time.Time       `json:"created_at,omitempty"`
	}
)
