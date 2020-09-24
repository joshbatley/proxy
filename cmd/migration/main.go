package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./storage.db")
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"ql", driver)
	if err != nil {
		log.Fatalf("HERE --- %v", err)
	}
	err = m.Up()
	fmt.Println(err)
}
