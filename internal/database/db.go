package database

import (
	"log"

	"github.com/jmoiron/sqlx"
)

// Conn create, Open set up DB
func Conn() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", "./storage.db")
	if err != nil {
		log.Panic(err)
	}

	return db
}
