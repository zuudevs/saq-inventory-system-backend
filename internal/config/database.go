package config

import (
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func NewDatabase(path string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite", path)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
