package config

import (
    "fmt"

    "github.com/jmoiron/sqlx"
    _ "github.com/go-sql-driver/mysql"
)

func NewDatabase(
    host,
    port,
    user,
    pass,
    name string,
) (*sqlx.DB, error) {

    dsn := fmt.Sprintf(
        "%s:%s@tcp(%s:%s)/%s?parseTime=true",
        user,
        pass,
        host,
        port,
        name,
    )

    db, err := sqlx.Connect("mysql", dsn)
    if err != nil {
        return nil, err
    }

    return db, nil
}