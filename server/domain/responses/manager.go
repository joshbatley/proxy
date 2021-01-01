package responses

import "github.com/google/uuid"

// Response returns struct from the database
type Response struct {
	ID         uuid.UUID `db:"ID"`
	EndpointID string    `db:"EndpointID"`
	Status     int       `db:"Status"`
	URL        string    `db:"URL"`
	Method     string    `db:"Method"`
	// Returns Headers as 'foo=bar; baz, other \n'
	Headers  string `db:"Headers"`
	Body     []byte `db:"Body"`
	DateTime int64  `db:"DateTime"`
}

// Repository -
type Repository interface {
	Get(url string, endpoint uuid.UUID, method string, status int) (*Response, error)
	ListByEndpoint(endpoint uuid.UUID, limit int, skip int) ([]Response, error)
	Save(id uuid.UUID, url string, h string, b []byte, st int, m string, e uuid.UUID) error
	Delete(id uuid.UUID) error
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
func (m *Manager) Get(url string, endpoint uuid.UUID, method string, status int) (*Response, error) {
	return m.repo.Get(url, endpoint, method, status)
}

// ListByEndpoint -
func (m *Manager) ListByEndpoint(endpoint uuid.UUID, limit int, skip int) ([]Response, error) {
	return m.repo.ListByEndpoint(endpoint, limit, skip)
}

// Save -
func (m *Manager) Save(
	id uuid.UUID, url string, head string, body []byte, status int, method string, endpointID uuid.UUID,
) error {
	if id == uuid.Nil {
		id = uuid.New()
	}

	return m.repo.Save(id, url, head, body, status, method, endpointID)
}

// Delete -
func (m *Manager) Delete(id uuid.UUID) error {
	return m.repo.Delete(id)
}
