package rules

import "github.com/jmoiron/sqlx"

type sqlRepo struct {
	db *sqlx.DB
}

// NewSQLRepository create new repository
func NewSQLRepository(db *sqlx.DB) *sqlRepo {
	return &sqlRepo{
		db: db,
	}
}

func (r *sqlRepo) Get(id int64) ([]Rule, error) {
	arr := []Rule{}
	err := r.db.Select(&arr, `
		SELECT Pattern, SaveResponse, ForceCors, Expiry, SkipOffline
		FROM Rules
		WHERE CollectionID=?
	`, id)

	if err != nil {
		return nil, err
	}

	return arr, nil
}
