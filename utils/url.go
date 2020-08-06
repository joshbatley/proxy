package utils

import (
	"net/url"
	"regexp"
	"strings"
)

// FormatURL -
func FormatURL(u string) *url.URL {
	r := regexp.MustCompile(`(?:/query\?q=)(.{0,})`)
	s := string(r.ReplaceAll([]byte(u), []byte("$1")))
	if !strings.Contains(s, "?") {
		s = strings.Replace(s, "&", "?", 1)
	}
	formattedURL, err := url.Parse(s)

	if err != nil {
		panic(err)
	}

	return formattedURL
}
