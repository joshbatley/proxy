package responses

import "database/sql"

// Response returns struct from the database
type Response struct {
	ID     int    `db:"ID"`
	Status int    `db:"Status"`
	URL    string `db:"URL"`
	// Returns Headers as 'foo=bar; baz, other \n'
	Headers  string `db:"Headers"`
	Body     []byte `db:"Body"`
	DateTime int64  `db:"DateTime"`
}

type repository interface {
	Get(u string, col int64, method string) (*Response, error)
	GetAllByCol(col int64) (*[]Response, error)
	Save(url string, h string, b []byte, st int, m string, e int64) error
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

func (m *Manager) Get(u string, col int64, method string) (*Response, error) {
	r, err := m.repo.Get(u, col, method)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return r, nil
}

func (m *Manager) GetAllByCol(col int64) (*[]Response, error) {
	r, err := m.repo.GetAllByCol(col)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return r, nil
}

func (m *Manager) Save(
	url string, head string, body []byte, status int, method string, endpointID int64,
) error {
	return m.repo.Save(url, head, body, status, method, endpointID)
}
