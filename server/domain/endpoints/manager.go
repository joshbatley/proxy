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

type repository interface {
	Get(url string, method string, col int64) (*Endpoint, error)
	GetByCollectionID(id int64) (*[]Endpoint, error)
	GetByID(id uuid.UUID) (*Endpoint, error)
	Save(url string, method string, col int64) (uuid.UUID, error)
	List(limit int, skip int) ([]Endpoint, error)
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

// List returns paginated response for endpoints
func (m *Manager) List(limit int, skip int) ([]Endpoint, error) {
	return m.repo.List(limit, skip)
}

// Get return all response by url, collection, method
func (m *Manager) Get(url string, method string, col int64) (*Endpoint, error) {
	return m.repo.Get(url, method, col)
}

// GetByCollectionID return all response by collection ID
func (m *Manager) GetByCollectionID(id int64) (*[]Endpoint, error) {
	return m.repo.GetByCollectionID(id)
}

// GetByID return all response by ID
func (m *Manager) GetByID(id uuid.UUID) (*Endpoint, error) {
	return m.repo.GetByID(id)
}

// Save new endpoint
func (m *Manager) Save(url string, method string, col int64) (uuid.UUID, error) {
	return m.repo.Save(url, method, col)
}
