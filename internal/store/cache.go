package store

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/joshbatley/proxy"
)

// CacheRepository setup repository with database
type CacheRepository struct {
	Database *sqlx.DB
}

// GetCache check DB for cached data based off url and collection
func (c *CacheRepository) GetCache(u string, col int64) (*proxy.CacheRow, error) {
	{
		err := c.Database.QueryRowx(`SELECT 1 FROM collection WHERE id=?`, col).Scan()
		if err == sql.ErrNoRows {
			log.Println("Collection not found")
			return nil, proxy.MissingColErr(err)
		}
	}

	var d proxy.CacheRow
	err := c.Database.QueryRowx(`
		SELECT id, body, status, headers, url
		FROM cache
		WHERE url=? AND collection=?
	`, u, col).StructScan(&d)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return &d, nil
}

// SaveCache saves the proxy request to the DB
func (c *CacheRepository) SaveCache(r *proxy.Record) error {
	_, err := c.Database.NamedExec(`
	INSERT INTO cache (
		url, headers, body, status, method, datetime, collection
	) VALUES (
		:url, :headers, :body, :status, :method, :datetime, :collection
	);`, &r)
	if err != nil {
		return err
	}
	return nil
}
