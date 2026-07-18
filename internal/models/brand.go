package models

import (
	"database/sql"
	"time"
)

type Brand struct {
	ID        uint64         `db:"id" json:"id"`
	Name      string         `db:"name" json:"name"`
	Slug      sql.NullString `db:"slug" json:"slug"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at"`
}
