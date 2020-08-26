package store

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/joshbatley/proxy"
)

// CacheStore setup repository with database
type CacheStore struct {
	Database *sqlx.DB
}

// GetAllCacheByColID returns all cache for collcetion
func (c *CacheStore) GetAllCacheByColID(col int64) (*[]proxy.CacheRow, error) {
	d := []proxy.CacheRow{}
	err := c.Database.Select(&d, `
		SELECT id, body, status, headers, url
		FROM cache
		WHERE collection=?
	`, col)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &d, nil
}

// GetCache return all cache where url and collection
func (c *CacheStore) GetCache(u string, col int64) (*proxy.CacheRow, error) {
	d := proxy.CacheRow{}
	err := c.Database.QueryRowx(`
		SELECT id, body, status, headers, url
		FROM cache
		WHERE url=? AND collection=?
	`, u, col).StructScan(&d)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &d, nil
}

// SaveCache saves the proxy request to the DB
func (c *CacheStore) SaveCache(r *proxy.Record) error {
	_, err := c.Database.NamedExec(`
	INSERT INTO cache (
		url, headers, body, status, method, datetime, collection
	) VALUES (
		:url, :headers, :body, :status, :method, :datetime, :collection
	);`, r)

	if err != nil {
		return err
	}
	return nil
}
