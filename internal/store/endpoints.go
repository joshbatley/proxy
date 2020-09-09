package store

import (
	"database/sql"
	"log"

	"github.com/joshbatley/proxy"
)

// GetOrAddEndpoint returns endpoint
func (s *Store) GetOrAddEndpoint(url string, method string, col int64) (*proxy.Endpoint, error) {
	e, err := s.getEndpoint(url, method, col)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if e != nil {
		return e, nil
	}

	id, err := s.AddEndpoint(
		col,
		url,
		method,
	)

	e, err = s.getEndpointByID(id)
	if err != nil {
		return nil, err
	}

	return e, nil
}

// AddEndpoint adds a endpoint
func (s *Store) AddEndpoint(col int64, url string, method string) (int64, error) {
	r, err := s.Database.Exec(`
	INSERT INTO Endpoints (
		CollectionID, URL, Method, PreferedStatus
	) VALUES (?, ?, ?, 200);`,
		col, url, method,
	)

	id, err := r.LastInsertId()
	log.Println(id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Store) getEndpoint(url string, method string, col int64) (*proxy.Endpoint, error) {
	e := proxy.Endpoint{}
	err := s.Database.QueryRowx(`
		SELECT
			ID, PreferedStatus, Method, URL
		FROM Endpoints WHERE URL=? AND Method=? and CollectionId=?
	`, url, method, col).StructScan(&e)

	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (s *Store) getEndpointByID(id int64) (*proxy.Endpoint, error) {
	e := proxy.Endpoint{}
	err := s.Database.QueryRowx(`
		SELECT ID, PreferedStatus, Method, URL FROM Endpoints WHERE id=?
	`, id).StructScan(&e)

	if err != nil {
		return nil, err
	}

	return &e, nil
}
