package store

import (
	"database/sql"

	"github.com/joshbatley/proxy"
)

// GetCollections get all the collections
func (s *Store) GetCollections() (*[]proxy.Collection, error) {
	cols := []proxy.Collection{}
	err := s.Database.Select(&cols, `
		SELECT * FROM Collections
	`)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return &cols, nil
}

// GetCollection by id
func (s *Store) GetCollection(id int64) (*proxy.Collection, error) {
	col := proxy.Collection{}
	err := s.Database.QueryRowx(`
		SELECT * FROM Collections WHERE ID=?
	`, id).StructScan(&col)

	if err != nil {
		return nil, err
	}
	return &col, nil
}

// SaveCollection add new collection
func (s *Store) SaveCollection(name string) (*proxy.Collection, error) {
	d, err := s.Database.NamedExec(`
		INSERT INTO Collections (
			Name
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
