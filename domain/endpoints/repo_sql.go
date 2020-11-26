package endpoints

import (
	"database/sql"

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

// List -
func (r *SQLRepo) List(limit int, skip int) ([]Endpoint, error) {
	cols := []Endpoint{}
	err := r.db.Select(&cols, `
		SELECT ID, PreferedStatus, Method, URL FROM Endpoints LIMIT ? OFFSET ?
	`, limit, skip)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return cols, nil
}

// Get -
func (r *SQLRepo) Get(url string, method string, col int64) (*Endpoint, error) {
	e := Endpoint{}
	err := r.db.QueryRowx(`
		SELECT
			ID, PreferedStatus, Method, URL
		FROM Endpoints WHERE URL=? AND Method=? and CollectionId=?
	`, url, method, col).StructScan(&e)

	if err == sql.ErrNoRows {
		return nil, fail.ErrNoData
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &e, nil
}

// GetByID -
func (r *SQLRepo) GetByID(id int64) (*[]Endpoint, error) {
	e := []Endpoint{}
	err := r.db.Select(&e, `
		SELECT ID, PreferedStatus, Method, URL FROM Endpoints WHERE CollectionID=?
	`, id)

	if err != nil {
		return nil, err
	}

	return &e, nil
}

// Save -
func (r *SQLRepo) Save(url string, method string, col int64) (uuid.UUID, error) {
	id := uuid.New()
	_, err := r.db.Exec(`
	INSERT INTO Endpoints (
		id, CollectionID, URL, Method, PreferedStatus
	) VALUES (?, ?, ?, ?, 200);`,
		id, col, url, method,
	)

	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
