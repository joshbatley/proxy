package store

import "github.com/jmoiron/sqlx"

// RulesStore setup repository with database
type RulesStore struct {
	Database *sqlx.DB
}
