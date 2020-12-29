package responses

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joshbatley/proxy/internal/fail"
)

// SQLRepo -
type SQLRepo struct {
	db *sqlx.DB
}

// NewSQLRepository create new repository
func NewSQLRepository(db *sqlx.DB) *SQLRepo {
	return &SQLRepo{
		db: db,
	}
}

// ListByEndpoint returns all response for Collection
func (r *SQLRepo) ListByEndpoint(endpoint string, limit int, skip int) ([]Response, error) {
	d := []Response{}
	err := r.db.Select(&d, `
		SELECT *
		FROM Responses
		WHERE EndpointID=?
		LIMIT ? OFFSET ?
	`, endpoint, limit, skip)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return d, nil
}

// Get return all response where url and collection
func (r *SQLRepo) Get(url string, endpoint uuid.UUID, method string, status int) (*Response, error) {
	d := Response{}
	err := r.db.QueryRowx(`
		SELECT ID, Body, Status, Headers, URL, DateTime
		FROM Responses
		WHERE URL=? AND EndpointId=? AND Method=? AND Status=?
	`, url, endpoint, method, status).StructScan(&d)

	if err == sql.ErrNoRows {
		return nil, fail.ErrNoData
	}

	if err != nil {
		return nil, err
	}

	return &d, nil
}

// Save saves the proxy request to the DB
func (r *SQLRepo) Save(id uuid.UUID, url string, h string, b []byte, st int, m string, e uuid.UUID) error {
	_, err := r.db.NamedExec(`
	INSERT OR REPLACE INTO Responses (
		ID, URL, Headers, Body, Status, Method, DateTime, EndpointID
	) VALUES (
		:id, :url, :headers, :body, :status, :method, :datetime, :endpoint
	);`, map[string]interface{}{
		"id":       id,
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

// Delete -
func (r *SQLRepo) Delete(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM Responses WHERE id=?", id)

	if err != nil {
		return err
	}

	return nil
}
