package utils

import (
	"net/url"
	"regexp"
	"strings"
)

// FormatURL takes the whole url and get only the requested
func FormatURL(u string) *url.URL {
	r := regexp.MustCompile(`(?:/query\?q=)(.{0,})`)
	s := string(r.ReplaceAll([]byte(u), []byte("$1")))
	// Fix for when formatting removes the extra "?"
	if !strings.Contains(s, "?") {
		s = strings.Replace(s, "&", "?", 1)
	}

	formattedURL, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return formattedURL
}
