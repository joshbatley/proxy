package responses

import "github.com/google/uuid"

// Response returns struct from the database
type Response struct {
	ID     string `db:"ID"`
	Status int    `db:"Status"`
	URL    string `db:"URL"`
	// Returns Headers as 'foo=bar; baz, other \n'
	Headers  string `db:"Headers"`
	Body     []byte `db:"Body"`
	DateTime int64  `db:"DateTime"`
}

// Repository -
type Repository interface {
	Get(u string, endpoint string, method string) (*Response, error)
	GetAllByCol(col int64) (*[]Response, error)
	Save(id string, url string, h string, b []byte, st int, m string, e uuid.UUID) error
	Delete(id string) error
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
func (m *Manager) Get(u string, endpoint string, method string) (*Response, error) {
	return m.repo.Get(u, endpoint, method)
}

// GetAllByCol -
func (m *Manager) GetAllByCol(col int64) (*[]Response, error) {
	return m.repo.GetAllByCol(col)
}

// Save -
func (m *Manager) Save(
	id string, url string, head string, body []byte, status int, method string, endpointID uuid.UUID,
) error {
	if len(id) == 0 {
		id = uuid.New().String()
	}

	return m.repo.Save(id, url, head, body, status, method, endpointID)
}

// Delete -
func (m *Manager) Delete(id string) error {
	return m.repo.Delete(id)
}
