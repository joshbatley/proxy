package responses

import (
	"database/sql"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type sqlRepo struct {
	db *sqlx.DB
}

// NewSQLRepository create new repository
func NewSQLRepository(db *sqlx.DB) *sqlRepo {
	return &sqlRepo{
		db: db,
	}
}

// GetAllResponseByColID returns all response for collcetion
func (r *sqlRepo) GetAllByCol(col int64) (*[]Response, error) {
	d := []Response{}
	err := r.db.Select(&d, `
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
func (r *sqlRepo) Get(u string, col int64, method string) (*Response, error) {
	d := Response{}
	err := r.db.QueryRowx(`
		SELECT ID, Body, Status, Headers, URL, DateTime
		FROM Responses
		WHERE URL=? AND EndpointId=? AND Method=?
	`, u, col, method).StructScan(&d)

	if err != nil && err != sql.ErrNoRows {
		log.Fatalln(err)
		return nil, err
	}

	return &d, nil
}

// SaveResponse saves the proxy request to the DB
func (r *sqlRepo) Save(url string, h string, b []byte, st int, m string, e int64) error {
	_, err := r.db.NamedExec(`
	INSERT INTO Responses (
		URL, Headers, Body, Status, Method, DateTime, EndpointID
	) VALUES (
		:url, :headers, :body, :status, :method, :datetime, :endpoint
	);`, map[string]interface{}{
		"url":      url,
		"headers":  h,
		"body":     b,
		"status":   st,
		"method":   m,
		"datetime": time.Now().Unix(),
		"endpoint": e,
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *sqlRepo) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM Responses WHERE id=?", id)

	if err != nil {
		return err
	}

	return nil
}
