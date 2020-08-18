package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/joshbatley/proxy/domain"
)

// CacheRepository -
type CacheRepository struct {
	Database *sqlx.DB
}

// CacheRow
type CacheRow struct {
	Status int
	URL    string
	// Returns Headers as 'foo=bar; baz, other \n'
	Headers string
	Body    []byte
}

const (
	selectCacheSQL = `
	SELECT body, status, headers, url FROM cache WHERE url=? AND collection =?
	`
	insertCacheSQL = `
	INSERT INTO cache (
		url, headers, body, status, method, datetime, collection
	) VALUES (
		:url, :headers, :body, :status, :method, :datetime, :collection
	);`
)

// GetCache -
func (c *CacheRepository) GetCache(u string, col int64) (*CacheRow, error) {
	// CHECK COLLECTION
	tx := c.Database.MustBegin()
	row := tx.QueryRowx(selectCacheSQL, u, col)
	var d CacheRow

	err := row.StructScan(&d)
	tx.Commit()
	return &d, err
}

// SaveCache -
func (c *CacheRepository) SaveCache(r *domain.Record) error {
	tx := c.Database.MustBegin()

	_, err := tx.NamedExec(insertCacheSQL, &r)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
