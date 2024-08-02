package db

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func DBConnection() (*sql.DB, error) {
	driver := "pgx"
	connstring := "postgres://postgres:didi@localhost:5432/event"
	db, err := sql.Open(driver, connstring)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
