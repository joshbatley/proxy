package collections

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/joshbatley/proxy/internal/fail"
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

func (r *sqlRepo) List() (*[]Collection, error) {
	cols := []Collection{}
	err := r.db.Select(&cols, `
		SELECT * FROM Collections
	`)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return &cols, nil
}

func (r *sqlRepo) Get(id int64) (*Collection, error) {
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

func (r *sqlRepo) Save(name string) (*Collection, error) {
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
