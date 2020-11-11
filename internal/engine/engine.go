package engine

import (
	"net/url"
	"regexp"
	"time"

	"github.com/joshbatley/proxy/internal/connection"
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
	url         *url.URL
	collection  int64
	rules       []Rule
	matchedRule *Rule
}

// Rule a single rule
type Rule struct {
	Pattern      string
	SaveResponse int
	ForceCors    int
	Expiry       int
	SkipOffline  int
}

// LoadRules pass in the request params and gets the rules
func (r *RuleEngine) LoadRules(url *url.URL, c int64, rules []Rule) {
	r.url = url
	r.collection = c
	r.rules = rules
}

// EnableCors Check rules to see if cors are enabled
func (r *RuleEngine) EnableCors() bool {
	rule := r.checkRules()
	return rule.ForceCors == 1
}

// CheckStore checks if store would be saved by using SaveResponse
func (r *RuleEngine) CheckStore() bool {
	rule := r.checkRules()
	return rule.SaveResponse == 1
}

// HasExpired Check rules to see expiry yime
func (r *RuleEngine) HasExpired(d int64) bool {
	rule := r.checkRules()
	exp := time.Unix(d, 0).Add(time.Second * time.Duration(rule.Expiry))
	if rule.SkipOffline == 1 {
		return exp.Before(time.Now())
	}
	if !connection.IsOnline(nil) {
		return false
	}
	return exp.Before(time.Now())
}

func (r *RuleEngine) checkRules() *Rule {
	if r.matchedRule == nil {
		for _, i := range r.rules {
			temp := regexp.MustCompilePOSIX(i.Pattern)
			matched := temp.Match([]byte(r.url.String()))
			if matched {
				r.matchedRule = &i
			}
		}
	}
	return r.matchedRule
}
