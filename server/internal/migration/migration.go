package migration

import (
	"database/sql"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

// StartUp Applies migration to sqlite DB
func StartUp() error {
	db, err := sql.Open("sqlite3", "./storage.db")
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	p, err := filepath.Abs("../migrations")
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+p,
		"ql", driver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
