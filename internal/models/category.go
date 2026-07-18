package models

import (
	"database/sql"
	"time"
)

type Category struct {
	ID          uint64         `db:"id" json:"id"`
	Name        string         `db:"name" json:"name"`
	Slug        string         `db:"slug" json:"slug"`
	Description sql.NullString `db:"description" json:"description"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at" json:"updated_at"`
}
