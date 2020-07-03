package db

import (
	"database/sql"
	"log"
	"os"
)

// Conn -
func Conn() {
	os.Remove("./foo.db")

	db, err := sql.Open("sqlite3", "./foo.db")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

// func initIfNew
// func loadExisting
