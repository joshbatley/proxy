package store

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/joshbatley/proxy"
)

// CollectionStore setup repository with database
type CollectionStore struct {
	Database *sqlx.DB
}

// GetCollections get all the collections
func (c *CollectionStore) GetCollections() (*[]proxy.Collection, error) {
	cols := []proxy.Collection{}
	err := c.Database.Select(&cols, `
		SELECT * FROM collection
	`)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return &cols, nil
}

// GetCollection by id
func (c *CollectionStore) GetCollection(id int64) (*proxy.Collection, error) {
	col := proxy.Collection{}
	err := c.Database.QueryRowx(`
		SELECT * FROM collection WHERE id=?
	`, id).StructScan(&col)

	if err != nil {
		return nil, err
	}
	return &col, nil
}

// SaveCollection add new collection
func (c *CollectionStore) SaveCollection(name string) (*proxy.Collection, error) {
	d, err := c.Database.NamedExec(`
		INSERT INTO collection (
			name
		) VALUES (
			:name
		)
	`, name)
	if err != nil {
		return nil, err
	}

	id, _ := d.LastInsertId()

	return &proxy.Collection{
		ID:   id,
		Name: name,
	}, nil
}
