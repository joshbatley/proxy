package engine

import (
	"log"
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
	url             *url.URL
	collection      int64
	rules           []Rule
	matchedRule     *Rule
	healthCheckURLS []string
}

// Rule a single rule
type Rule struct {
	Pattern      string
	SaveResponse int
	ForceCors    int
	Expiry       int
	SkipOffline  int
	Delay        int
	RemapRegex   string
}

// LoadRules pass in the request params and gets the rules
func (r *RuleEngine) LoadRules(url *url.URL, c int64, rules []Rule, healthCheckURLS []string) {
	r.url = url
	r.collection = c
	r.rules = rules
	r.healthCheckURLS = healthCheckURLS
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
	if connection.IsOffline(r.healthCheckURLS) {
		log.Println("connection is offline")
		return false
	}
	return exp.Before(time.Now())
}

// Remapper remaps the urls
func (r *RuleEngine) Remapper() *url.URL {
	rule := r.checkRules()
	if len(rule.RemapRegex) != 0 {
		temp := regexp.MustCompile(rule.Pattern)
		s := temp.ReplaceAllString(r.url.String(), rule.RemapRegex)
		newURL, err := url.ParseRequestURI(s)
		if err == nil {
			return newURL
		}
	}
	return r.url
}

// GetSleepTime Get the current rules delay time
func (r *RuleEngine) GetSleepTime() int64 {
	rule := r.checkRules()
	return int64(rule.Delay)
}

func (r *RuleEngine) checkRules() *Rule {
	if r.matchedRule == nil {
		for _, i := range r.rules {
			temp := regexp.MustCompile(i.Pattern)
			matched := temp.Match([]byte(r.url.String()))
			if matched {
				r.matchedRule = &i
			}
		}
	}
	return r.matchedRule
}
