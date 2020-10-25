package database

import (
	"github.com/jmoiron/sqlx"
)

// Conn create, Open set up DB
func Conn() (*sqlx.DB, error) {
	return sqlx.Open("sqlite3", "./storage.db")
}
