package engine

import (
	"database/sql"
	"regexp"
	"time"

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
	StateProxyAndSave
	StateReturnRefresh
	StateOffline
	StateUpdate
)

// RuleEngine powers all rules
type RuleEngine struct {
	store       *store.Store
	params      *utils.Params
	rules       []proxy.Rule
	matchedRule *proxy.Rule
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

func (r *RuleEngine) reset(p *utils.Params) {
	if r.params != nil {
		if r.params.QueryURL != p.QueryURL && r.params.Collection != p.Collection {
			r.matchedRule = nil
			r.rules = make([]proxy.Rule, 0)
		}
	}
}

// LoadRules pass in the request params and gets the rules
func (r *RuleEngine) LoadRules(p *utils.Params) error {
	r.reset(p)
	_, err := r.store.GetCollection(p.Collection)
	if err == sql.ErrNoRows {
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
func (r *RuleEngine) GetState() State {
	rule := r.checkRules()
	if rule.SaveResponse == 1 {
		return StateSaving
	}
	return StateProxy
}

// EnableCors Check rules to see if cors are enabled
func (r *RuleEngine) EnableCors() bool {
	rule := r.checkRules()
	return rule.ForceCors == 1
}

// HasExpired Check rules to see expiry yime
func (r *RuleEngine) HasExpired(d int64) bool {
	rule := r.checkRules()
	exp := time.Unix(d, 0).Add(time.Second * time.Duration(rule.Expiry))
	return exp.Before(time.Now())
}

func (r *RuleEngine) checkRules() *proxy.Rule {
	if r.matchedRule == nil {
		for _, i := range r.rules {
			temp := regexp.MustCompilePOSIX(i.Pattern)
			matched := temp.Match([]byte(r.params.QueryURL.String()))
			if matched {
				r.matchedRule = &i
			}
		}
	}
	return r.matchedRule
}
