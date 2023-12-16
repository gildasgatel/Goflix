package db

import (
	"database/sql"

	"goflix/config"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	Setup() error
	Close()
}

type DbSqlite struct {
	sqlite *sql.DB
}

func New() Storage {
	return &DbSqlite{}
}

func (db *DbSqlite) Setup() error {
	var err error
	db.sqlite, err = sql.Open(config.DRIVE_NAME, config.DATA_SOURCE_NAME)
	if err != nil {
		return err
	}
	return nil
}

func (db *DbSqlite) Close() {
	db.sqlite.Close()
}
