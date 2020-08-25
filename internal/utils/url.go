package utils

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/joshbatley/proxy"
)

// Params for a request
type Params struct {
	QueryURL   *url.URL
	Collection int64
}

// ParseParams takes the url and returns Params
func ParseParams(ou map[string]string, q *url.URL) (*Params, error) {
	u := q.String()

	if strings.HasPrefix(u, "/") || strings.HasPrefix(u, "/"+ou["collection"]+"/") {
		u = strings.TrimPrefix(u, "/"+ou["collection"]+"/")
		u = strings.TrimPrefix(u, "/")
	}

	formattedURL, err := url.ParseRequestURI(u)
	if err != nil {
		return &Params{}, proxy.URLInvalidErr(err)
	}

	var c int64
	c, err = strconv.ParseInt(ou["collection"], 0, 64)
	if err != nil {
		c = 1
	}

	return &Params{
		QueryURL:   formattedURL,
		Collection: c,
	}, nil
}
