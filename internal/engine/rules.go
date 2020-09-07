package engine

import (
	"net/url"
	"regexp"

	"github.com/joshbatley/proxy"
)

// State Rule current state
type State int

// Possible States
const (
	StateSaving State = iota + 1
	StateProxy
)

// Engine -
func Engine(rs []proxy.Rule, URL *url.URL) State {
	for _, r := range rs {
		temp := regexp.MustCompilePOSIX(r.Pattern)
		matched := temp.Match([]byte(URL.String()))
		if matched && r.Cache == 1 {
			return StateSaving
		}
	}

	return StateProxy
}
