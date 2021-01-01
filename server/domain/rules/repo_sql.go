package rules

import "github.com/jmoiron/sqlx"

// SQLRepo requires DB
type SQLRepo struct {
	db *sqlx.DB
}

// NewSQLRepository create new repository
func NewSQLRepository(db *sqlx.DB) *SQLRepo {
	return &SQLRepo{
		db: db,
	}
}

// GetByCollectionID get all rules by the colleciton ID
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

// ListByCollectionID paginatied request by the colleciton ID
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
