package endpoints

import (
	"github.com/google/uuid"
)

// Endpoint returns a single endpoint
type Endpoint struct {
	ID     string `db:"ID"`
	Status int    `db:"PreferedStatus"`
	Method string `db:"Method"`
	URL    string `db:"URL"`
}

// Repository -
type Repository interface {
	Get(url string, method string, col int64) (*Endpoint, error)
	GetByID(id int64) (*Endpoint, error)
	Save(col int64, url string, method string) (uuid.UUID, error)
}

// Manager -
type Manager struct {
	repo Repository
}

//NewManager create new manager
func NewManager(r Repository) *Manager {
	return &Manager{
		repo: r,
	}
}

// Get -
func (m *Manager) Get(url string, method string, col int64) (*Endpoint, error) {
	return m.repo.Get(url, method, col)
}

// GetByID -
func (m *Manager) GetByID(id int64) (*Endpoint, error) {
	return m.repo.GetByID(id)
}

// Save -
func (m *Manager) Save(col int64, url string, method string) (uuid.UUID, error) {
	return m.repo.Save(col, url, method)
}
