package collections

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/joshbatley/proxy/server/internal/fail"
)

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
func (r *SQLRepo) List(limit int, skip int) ([]Collection, error) {
	cols := []Collection{}
	err := r.db.Select(&cols, `
		SELECT * FROM Collections LIMIT ? OFFSET ?
	`, limit, skip)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return cols, nil
}

// Get -
func (r *SQLRepo) Get(id int64) (*Collection, error) {
	col := Collection{}
	err := r.db.QueryRowx(`
		SELECT * FROM Collections WHERE ID=?
	`, id).StructScan(&col)

	if err == sql.ErrNoRows {
		return nil, fail.ErrMissingCol
	}

	if err != nil {
		return nil, err
	}
	return &col, nil
}

// Save -
func (r *SQLRepo) Save(name string) (*Collection, error) {
	d, err := r.db.NamedExec(`
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

	return &Collection{
		ID:   id,
		Name: name,
	}, nil
}
