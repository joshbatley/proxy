package database

import (
	"log"

	"github.com/jmoiron/sqlx"
)

const (
	collectionTableSQL = `
	CREATE TABLE IF NOT EXISTS collection(
		id INTEGER NOT NULL PRIMARY KEY,
		friendlyname TEXT NOT NULL,
		UNIQUE(friendlyname)
	);
	INSERT OR IGNORE INTO collection (friendlyname) VALUES ("DEFAULT")
	`
	cacheTableSQL = `
	CREATE TABLE IF NOT EXISTS cache(
		id INTEGER NOT NULL PRIMARY KEY,
		collection INTEGER NOT NULL,
		url TEXT NOT NULL,
		headers TEXT,
		body BLOB,
		status INTEGER,
		method TEXT,
		datetime INTEGER,
		FOREIGN KEY(collection) REFERENCES collection(id)
	);`
	rulesTableSQL = `
	CREATE TABLE IF NOT EXISTS rules(
		id INTEGER NOT NULL PRIMARY KEY,
		collection INTEGER NOT NULL,
		pattern TEXT NOT NULL,
		cache INTEGER,
		expiry INTEGER,
		offlinecache INTEGER,
		FOREIGN KEY(collection) REFERENCES collection(id)
	);`
)

var db *sqlx.DB

// Conn -
func Conn() *sqlx.DB {
	// os.Remove("./storage.db")

	conn, err := sqlx.Open("sqlite3", "./storage.db")
	if err != nil {
		log.Fatal(err)
	}

	if db == nil {
		db = conn
	}

	setup()

	return db
}

func setup() {
	_, err := db.Exec(collectionTableSQL)
	if err != nil {
		log.Panic(err)
	}

	_, err = db.Exec(cacheTableSQL)
	if err != nil {
		log.Panic(err)
	}

	_, err = db.Exec(rulesTableSQL)
	if err != nil {
		log.Panic(err)
	}
}
