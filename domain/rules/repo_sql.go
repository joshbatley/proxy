package rules

import "github.com/jmoiron/sqlx"

type SQLRepo struct {
	db *sqlx.DB
}

// NewSQLRepository create new repository
func NewSQLRepository(db *sqlx.DB) *SQLRepo {
	return &SQLRepo{
		db: db,
	}
}

func (r *SQLRepo) GetByCollectionID(id int64) ([]Rule, error) {
	arr := []Rule{}
	err := r.db.Select(&arr, `
		SELECT Pattern,
		SaveResponse,
		ForceCors,
		Expiry,
		SkipOffline,
		DelayTime,
		RemapRegex
		FROM Rules
		WHERE CollectionID=?
	`, id)

	if err != nil {
		return nil, err
	}

	return arr, nil
}

func (r *SQLRepo) ListByCollectionID(collection string, limit int, skip int) ([]Rule, error) {
	rules := []Rule{}
	err := r.db.Select(&rules, `
		SELECT Pattern,
		SaveResponse,
		ForceCors,
		Expiry,
		SkipOffline,
		DelayTime,
		RemapRegex
		FROM Rules
		WHERE CollectionID=?
		LIMIT ? OFFSET ?
	`, collection, limit, skip)

	if err != nil {
		return nil, err
	}

	return rules, nil
}
