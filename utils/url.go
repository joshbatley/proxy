package utils

import (
	"fmt"
	"net/url"
	"strings"
)

// Params for a request
type Params struct {
	QueryURL   *url.URL
	Collection string
}

type queryParseError struct {
	Message string
	Log     string
}

func (e *queryParseError) Error() string {
	return fmt.Sprintf("%s - %s", e.Message, e.Log)
}

// FormatURL takes the url and retursn Params
func FormatURL(ou map[string]string, q *url.URL) (*Params, error) {
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

	return &Params{
		QueryURL:   formattedURL,
		Collection: ou["collection"],
	}, nil
}
