package repository

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/joshbatley/proxy/domain"
	"github.com/joshbatley/proxy/utils"
)

// CacheRepository setup repository with database
type CacheRepository struct {
	Database *sqlx.DB
}

const (
	selectCacheSQL = `
	SELECT id, body, status, headers, url FROM cache WHERE url=? AND collection=?
	`
	selectCollectionSQL = `
	SELECT 1 FROM collection WHERE id=?
	`
	insertCacheSQL = `
	INSERT INTO cache (
		url, headers, body, status, method, datetime, collection
	) VALUES (
		:url, :headers, :body, :status, :method, :datetime, :collection
	);`
)

// GetCache check DB for cached data based off url and collection
func (c *CacheRepository) GetCache(u string, col int64) (*domain.CacheRow, error) {
	tx, err := c.Database.Beginx()

	{
		err := tx.QueryRowx(selectCollectionSQL, col).Scan()
		defer tx.Commit()

		if err == sql.ErrNoRows {
			log.Println("Collection not found")
			return nil, utils.ColMissingErr(err)
		}
	}

	var d domain.CacheRow
	err = tx.QueryRowx(selectCacheSQL, u, col).StructScan(&d)
	defer tx.Commit()

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return &d, nil
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
