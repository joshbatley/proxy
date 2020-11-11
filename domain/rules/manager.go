package rules

// Rule returns a single rule
type Rule struct {
	Pattern      string `db:"Pattern"`
	SaveResponse int    `db:"SaveResponse"`
	ForceCors    int    `db:"ForceCors"`
	Expiry       int    `db:"Expiry"`
	SkipOffline  int    `db:"SkipOffline"`
}

type repository interface {
	Get(id int64) ([]Rule, error)
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

// GetRules by collection ID
func (m *Manager) Get(id int64) ([]Rule, error) {
	return m.repo.Get(id)
}
