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

type repository interface {
	Get(url string, endpoint uuid.UUID, method string, status int) (*Response, error)
	ListByEndpoint(endpoint uuid.UUID, limit int, skip int) ([]Response, error)
	Save(id uuid.UUID, url string, h string, b []byte, st int, m string, e uuid.UUID) error
	Delete(id uuid.UUID) error
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

// Get return all response by url, endpoint, method and status
func (m *Manager) Get(url string, endpoint uuid.UUID, method string, status int) (*Response, error) {
	return m.repo.Get(url, endpoint, method, status)
}

// ListByEndpoint returns paginated response by endpoint
func (m *Manager) ListByEndpoint(endpoint uuid.UUID, limit int, skip int) ([]Response, error) {
	return m.repo.ListByEndpoint(endpoint, limit, skip)
}

// Save the proxy request to the DB
func (m *Manager) Save(
	id uuid.UUID, url string, head string, body []byte, status int, method string, endpointID uuid.UUID,
) error {
	if id == uuid.Nil {
		id = uuid.New()
	}

	return m.repo.Save(id, url, head, body, status, method, endpointID)
}

// Delete one required by ID
func (m *Manager) Delete(id uuid.UUID) error {
	return m.repo.Delete(id)
}
