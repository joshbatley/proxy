package store

import (
	"github.com/joshbatley/proxy"
)

// GetRules by collection ID
func (s *Store) GetRules(id int64) ([]proxy.Rule, error) {
	r := []proxy.Rule{}
	err := s.Database.Select(&r, `
		SELECT Pattern, SaveResponse, ForceCors, Expiry
		FROM Rules
		WHERE CollectionID=?
	`, id)

	if err != nil {
		return nil, proxy.InternalError(err)
	}

	return r, nil
}
