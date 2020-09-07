package store

import (
	"database/sql"

	"github.com/joshbatley/proxy"
)

// GetAllCacheByColID returns all cache for collcetion
func (s *Store) GetAllCacheByColID(col int64) (*[]proxy.CacheRow, error) {
	d := []proxy.CacheRow{}
	err := s.Database.Select(&d, `
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
func (s *Store) GetCache(u string, col int64) (*proxy.CacheRow, error) {
	d := proxy.CacheRow{}
	err := s.Database.QueryRowx(`
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
func (s *Store) SaveCache(r *proxy.Record) error {
	_, err := s.Database.NamedExec(`
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
