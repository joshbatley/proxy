package engine

import (
	"database/sql"
	"regexp"

	"github.com/joshbatley/proxy"
	"github.com/joshbatley/proxy/internal/store"
	"github.com/joshbatley/proxy/internal/utils"
)

// State Rule current state
type State int

// Possible States
const (
	StateSaving State = iota + 1
	StateProxy
)

// RuleEngine powers all rules
type RuleEngine struct {
	store  *store.Store
	params *utils.Params
	rules  []proxy.Rule
}

// NewEngine inits a new Engine
func NewEngine(store *store.Store) *RuleEngine {
	return &RuleEngine{
		store: store,
	}
}

func (r *RuleEngine) getRules() ([]proxy.Rule, error) {
	res, err := r.store.GetRules(r.params.Collection)
	if err != nil {
		return nil, proxy.InternalError(err)
	}
	return res, nil
}

// StartUp pass in the request params and gets the rules
func (r *RuleEngine) StartUp(p *utils.Params) error {
	if _, err := r.store.GetCollection(p.Collection); err == sql.ErrNoRows {
		return proxy.MissingColErr(err)
	}
	r.params = p

	rules, err := r.getRules()
	if err != nil {
		return err
	}

	r.rules = rules
	return nil
}

// GetState -
func (r *RuleEngine) GetState() (State, error) {
	for _, i := range r.rules {
		temp := regexp.MustCompilePOSIX(i.Pattern)
		matched := temp.Match([]byte(r.params.QueryURL.String()))
		if matched && i.Cache == 1 {
			return StateSaving, nil
		}
	}

	return StateProxy, nil
}
