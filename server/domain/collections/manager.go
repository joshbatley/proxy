package collections

import "database/sql"

// Collection returns struct from the database
type Collection struct {
	ID              int64          `db:"ID"`
	Name            string         `db:"Name"`
	HealthCheckURLs sql.NullString `db:"HealthCheckURLs"`
}

type repository interface {
	List(limit int, skip int) ([]Collection, error)
	Get(id int64) (*Collection, error)
	Save(name string) (*Collection, error)
}

// Manager requires repo
type Manager struct {
	repo repository
}

//NewManager create new manager
func NewManager(r repository) *Manager {
	return &Manager{
		repo: r,
	}
}

// List returns paginated response for Collection
func (m *Manager) List(limit int, skip int) ([]Collection, error) {
	return m.repo.List(limit, skip)
}

// Get collection by id
func (m *Manager) Get(id int64) (*Collection, error) {
	return m.repo.Get(id)
}

// Save new Collection
func (m *Manager) Save(name string) (*Collection, error) {
	return m.repo.Save(name)
}
