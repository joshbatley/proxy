package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/joshbatley/proxy/def"
)

const collectionTableSQL = `
	CREATE TABLE collections(
		id INTEGER NOT NULL PRIMARY KEY,
		name TEXT,
		baseurl TEXT
	);
`

const cacheTableSQL = `
	CREATE TABLE cache(
		id INTEGER NOT NULL PRIMARY KEY,
		collection INTEGER,
		url TEXT NOT NULL,
		headers TEXT,
		body BLOB,
		status INTEGER,
		method TEXT,
		dateTime INTEGER
	);
`

//		FOREIGN KEY(collection) REFERENCES collections(id)

const rulesTableSQL = `
	CREATE TABLE rules(
		id INTEGER NOT NULL PRIMARY KEY,
		collection INTEGER,
		pattern TEXT NOT NULL,
		cache INTEGER,
		expiry INTEGER,
		offlinecache INTEGER,
		FOREIGN KEY(collection) REFERENCES collections(id)
	);
`

var db *sql.DB

// Conn -
func Conn() {
	os.Remove("./storage.db")

	conn, err := sql.Open("sqlite3", "./storage.db")
	if err != nil {
		log.Fatal(err)
	}

	if db == nil {
		db = conn
	}

	setup()
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

func Query() {
	// db.Query(query string, args ...interface{})
	// row := db.QueryRow("SELECT url, body, headers, status, method FROM cache WHERE url = '?'", url.String())
	// err := row.Scan(&data.URL, &data.Body, &data.Headers, &data.Status, &data.Method)
	// if err == (sql.ErrNoRows) {
	// 	log.Println("no row")
	// 	return false
	// }
	// if err != nil {
	// 	log.Println(err)
	// 	return false
	// }
}

// Insert
func Insert(r def.Record) {
	tx, _ := db.Begin()
	stmt, _ := tx.Prepare(`INSERT INTO cache (url, headers, body, status, method) values (?,?,?,?,?)`)

	_, err := stmt.Exec(r.URLString(), r.HeadersToString(), r.Body, r.Status, r.Method)
	if err != nil {
		panic(err)
	}
	tx.Commit()
}
