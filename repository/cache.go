package repository

import (
	"database/sql"
	"log"

	"github.com/joshbatley/proxy/database"
	"github.com/joshbatley/proxy/domain"
)

// CacheRepository -
type CacheRepository struct {
	Database *sql.DB
}

// Cache -
type Cache struct {
	Status int
	URL    string
	// Returns Headers as 'foo=bar; baz, other \n'
	Header string
	Body   []byte
}

// GetCache -
func (c *CacheRepository) GetCache(u string) (Cache, error) {
	tx, _ := c.Database.Begin()
	log.Printf(u)
	row := tx.QueryRow("SELECT body, status, headers, url FROM cache WHERE url=?", u)
	var d Cache

	err := row.Scan(&d.Body, &d.Status, &d.Header, &d.URL)

	tx.Commit()
	return d, err
}

// SaveCache -
func (c *CacheRepository) SaveCache(r domain.Record) {
	database.Insert(r)
}
