package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/joshbatley/proxy/domain"
)

const cacheTableSQL = `
	CREATE TABLE IF NOT EXISTS cache(
		id INTEGER NOT NULL PRIMARY KEY,
		collection TEXT NOT NULL,
		url TEXT NOT NULL,
		headers TEXT,
		body BLOB,
		status INTEGER,
		method TEXT,
		datetime INTEGER
	);
`

//		FOREIGN KEY(collection) REFERENCES collections(id)

const rulesTableSQL = `
	CREATE TABLE IF NOT EXISTS rules(
		id INTEGER NOT NULL PRIMARY KEY,
		collection INTEGER,
		pattern TEXT NOT NULL,
		cache INTEGER,
		expiry INTEGER,
		offlinecache INTEGER,
		FOREIGN KEY(collection) REFERENCES cache(collection)
	);
`

var db *sql.DB

// Conn -
func Conn() *sql.DB {
	// os.Remove("./storage.db")

	conn, err := sql.Open("sqlite3", "./storage.db")
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
	_, err := db.Exec(cacheTableSQL)
	if err != nil {
		log.Panic(err)
	}

	_, err = db.Exec(rulesTableSQL)
	if err != nil {
		log.Panic(err)
	}
}

// Insert -
func Insert(r domain.Record) {
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare(`INSERT INTO cache (collection, url, headers, body, status, method, datetime) values (?,?,?,?,?,?, ?)`)

	_, err := stmt.Exec(r.URL.Host, r.URLString(), r.HeadersToString(), r.Body, r.Status, r.Method, time.Now())
	if err != nil {
		panic(err)
	}
	if err = tx.Commit(); err != nil {
		log.Panicln(err)
	} else {
		log.Println("Saving", r.URL.String())
	}
}
