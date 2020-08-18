package utils

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Params for a request
type Params struct {
	QueryURL   *url.URL
	Collection int64
}

type queryParseError struct {
	Message string
	Log     string
}

func (e *queryParseError) Error() string {
	return fmt.Sprintf("%s - %s", e.Message, e.Log)
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
		return &Params{}, &queryParseError{
			Message: "Requested URL is not valid",
			Log:     err.Error(),
		}
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
