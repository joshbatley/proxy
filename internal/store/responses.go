package store

import (
	"database/sql"

	"github.com/joshbatley/proxy"
)

// GetAllResponseByColID returns all response for collcetion
func (s *Store) GetAllResponseByColID(col int64) (*[]proxy.ResponseRow, error) {
	d := []proxy.ResponseRow{}
	err := s.Database.Select(&d, `
		SELECT ID, Body, Status, Headers, URL
		FROM Responses
		WHERE CollectionID=?
	`, col)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &d, nil
}

// GetResponse return all response where url and collection
func (s *Store) GetResponse(u string, col int64, method string) (*proxy.ResponseRow, error) {
	d := proxy.ResponseRow{}
	err := s.Database.QueryRowx(`
		SELECT ID, Body, Status, Headers, URL, DateTime
		FROM Responses
		WHERE URL=? AND EndpointId=? AND Method=?
	`, u, col, method).StructScan(&d)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &d, nil
}

// SaveResponse saves the proxy request to the DB
func (s *Store) SaveResponse(r *proxy.Response) error {
	_, err := s.Database.NamedExec(`
	INSERT INTO Responses (
		URL, Headers, Body, Status, Method, DateTime, EndpointID
	) VALUES (
		:url, :headers, :body, :status, :method, :datetime, :endpoint
	);`, r)

	if err != nil {
		return err
	}
	return nil
}
