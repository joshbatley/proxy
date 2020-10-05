package endpoints

import (
	"database/sql"

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

func (r *sqlRepo) Get(url string, method string, col int64) (*Endpoint, error) {
	e := Endpoint{}
	err := r.db.QueryRowx(`
		SELECT
			ID, PreferedStatus, Method, URL
		FROM Endpoints WHERE URL=? AND Method=? and CollectionId=?
	`, url, method, col).StructScan(&e)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &e, nil
}

func (r *sqlRepo) GetByID(id int64) (*Endpoint, error) {
	e := Endpoint{}
	err := r.db.QueryRowx(`
		SELECT ID, PreferedStatus, Method, URL FROM Endpoints WHERE id=?
	`, id).StructScan(&e)

	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *sqlRepo) Save(col int64, url string, method string) (int64, error) {
	row, err := r.db.Exec(`
	INSERT INTO Endpoints (
		CollectionID, URL, Method, PreferedStatus
	) VALUES (?, ?, ?, 200);`,
		col, url, method,
	)

	id, err := row.LastInsertId()

	if err != nil {
		return 0, err
	}

	return id, nil
}
