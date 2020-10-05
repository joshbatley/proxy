package endpoints

import (
	"log"
)

// Endpoint returns a single endpoint
type Endpoint struct {
	ID     int64  `db:"ID"`
	Status int    `db:"PreferedStatus"`
	Method string `db:"Method"`
	URL    string `db:"URL"`
}

type repository interface {
	Get(url string, method string, col int64) (*Endpoint, error)
	GetByID(id int64) (*Endpoint, error)
	Save(col int64, url string, method string) (int64, error)
}

type Manager struct {
	repo repository
}

//NewManager create new manager
func NewManager(r repository) *Manager {
	return &Manager{
		repo: r,
	}
}

func (m *Manager) GetOrSave(url string, method string, col int64) (*Endpoint, error) {
	endpoint, err := m.Get(url, method, col)
	if err != nil {
		log.Println(err)

		return nil, err
	}

	if endpoint != nil {
		return endpoint, nil
	}

	id, err := m.Save(col, url, method)

	endpoint, err = m.GetByID(id)

	if err != nil {
		return nil, err
	}

	return endpoint, nil
}

func (m *Manager) Get(url string, method string, col int64) (*Endpoint, error) {
	return m.repo.Get(url, method, col)
}

func (m *Manager) GetByID(id int64) (*Endpoint, error) {
	return m.repo.GetByID(id)
}

func (m *Manager) Save(col int64, url string, method string) (int64, error) {
	return m.repo.Save(col, url, method)
}
