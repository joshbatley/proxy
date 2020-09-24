package collections

import "database/sql"

// Collection returns struct from the database
type Collection struct {
	ID   int64  `db:"ID"`
	Name string `db:"Name"`
}

type repository interface {
	List() (*[]Collection, error)
	Get(id int64) (*Collection, error)
	Save(name string) (*Collection, error)
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

func (m *Manager) List() (*[]Collection, error) {
	c, err := m.repo.List()
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return c, err
}

func (m *Manager) Get(id int64) (*Collection, error) {
	return m.repo.Get(id)
}

func (m *Manager) Save(name string) (*Collection, error) {
	return m.repo.Save(name)
}
