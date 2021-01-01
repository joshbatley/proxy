package endpoints

import (
	"github.com/google/uuid"
)

// Endpoint returns a single endpoint
type Endpoint struct {
	ID           uuid.UUID `db:"ID"`
	Status       int       `db:"PreferedStatus"`
	Method       string    `db:"Method"`
	URL          string    `db:"URL"`
	CollectionID string    `db:"CollectionID"`
}

// Repository -
type Repository interface {
	Get(url string, method string, col int64) (*Endpoint, error)
	GetByColID(id int64) (*[]Endpoint, error)
	GetByID(id uuid.UUID) (*Endpoint, error)
	Save(url string, method string, col int64) (uuid.UUID, error)
	List(limit int, skip int) ([]Endpoint, error)
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

// List -
func (m *Manager) List(limit int, skip int) ([]Endpoint, error) {
	return m.repo.List(limit, skip)
}

// Get -
func (m *Manager) Get(url string, method string, col int64) (*Endpoint, error) {
	return m.repo.Get(url, method, col)
}

// GetByColID -
func (m *Manager) GetByColID(id int64) (*[]Endpoint, error) {
	return m.repo.GetByColID(id)
}

// GetByID -
func (m *Manager) GetByID(id uuid.UUID) (*Endpoint, error) {
	return m.repo.GetByID(id)
}

// Save -
func (m *Manager) Save(url string, method string, col int64) (uuid.UUID, error) {
	return m.repo.Save(url, method, col)
}
