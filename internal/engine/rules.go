package engine

import (
	"database/sql"
	"regexp"
	"time"

	"github.com/joshbatley/proxy/domain/collections"
	"github.com/joshbatley/proxy/domain/rules"
	"github.com/joshbatley/proxy/internal/fail"
	"github.com/joshbatley/proxy/internal/params"
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
	rules       *rules.Manager
	collection  *collections.Manager
	params      *params.Params
	r           []rules.Rule
	matchedRule *rules.Rule
}

// NewEngine inits a new Engine
func NewEngine(rules *rules.Manager, collection *collections.Manager) *RuleEngine {
	return &RuleEngine{
		rules:      rules,
		collection: collection,
	}
}

func (r *RuleEngine) getRules() ([]rules.Rule, error) {
	res, err := r.rules.Get(r.params.Collection)
	if err != nil {
		return nil, fail.InternalError(err)
	}
	return res, nil
}

func (r *RuleEngine) reset(p *params.Params) {
	if r.params != nil {
		if r.params.QueryURL != p.QueryURL && r.params.Collection != p.Collection {
			r.matchedRule = nil
			r.r = make([]rules.Rule, 0)
		}
	}
}

// LoadRules pass in the request params and gets the rules
func (r *RuleEngine) LoadRules(p *params.Params) error {
	r.reset(p)
	_, err := r.collection.Get(p.Collection)
	if err == sql.ErrNoRows {
		return fail.MissingColErr(err)
	}
	r.params = p

	rules, err := r.getRules()
	if err != nil {
		return err
	}

	r.r = rules
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

func (r *RuleEngine) checkRules() *rules.Rule {
	if r.matchedRule == nil {
		for _, i := range r.r {
			temp := regexp.MustCompilePOSIX(i.Pattern)
			matched := temp.Match([]byte(r.params.QueryURL.String()))
			if matched {
				r.matchedRule = &i
			}
		}
	}
	return r.matchedRule
}
