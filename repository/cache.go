package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/joshbatley/proxy/domain"
)

// CacheRepository setup repository with database
type CacheRepository struct {
	Database *sqlx.DB
}

// CacheRow returns struct from the database
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

// GetCache check DB for cached data based off url and collection
func (c *CacheRepository) GetCache(u string, col int64) (*CacheRow, error) {
	// CHECK COLLECTION
	tx := c.Database.MustBegin()
	row := tx.QueryRowx(selectCacheSQL, u, col)
	var d CacheRow

	err := row.StructScan(&d)
	tx.Commit()
	return &d, err
}

// SaveCache saves the proxy request to the DB
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
