package database

import (
	"io/ioutil"
	"log"

	"github.com/jmoiron/sqlx"
)

// Conn create, Open set up DB
func Conn() *sqlx.DB {
	// os.Remove("./storage.db")

	db, _ := sqlx.Open("sqlite3", "./storage.db")
	query, err := ioutil.ReadFile("./database/migrations/setup.sql")
	if err != nil {
		log.Panic(err)
	}
	if _, err := db.Exec(string(query)); err != nil {
		log.Panic(err)
	}

	return db
}
