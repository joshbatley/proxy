package rules

// Rule returns a single rule
type Rule struct {
	Pattern      string `db:"Pattern"`
	SaveResponse int    `db:"SaveResponse"`
	ForceCors    int    `db:"ForceCors"`
	Expiry       int    `db:"Expiry"`
	SkipOffline  int    `db:"SkipOffline"`
	Delay        int    `db:"DelayTime"`
	RemapRegex   string `db:"RemapRegex"`
}

type repository interface {
	GetByCollectionID(id int64) ([]Rule, error)
	ListByCollectionID(collection string, limit int, skip int) ([]Rule, error)
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

// GetByCollectionID by collection ID
func (m *Manager) GetByCollectionID(id int64) ([]Rule, error) {
	return m.repo.GetByCollectionID(id)
}

// ListByCollectionID -
func (m *Manager) ListByCollectionID(collection string, limit int, skip int) ([]Rule, error) {
	return m.repo.ListByCollectionID(collection, limit, skip)
}
