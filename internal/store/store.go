package store

import "github.com/jmoiron/sqlx"

// Store setup repository with database
type Store struct {
	Database *sqlx.DB
}
