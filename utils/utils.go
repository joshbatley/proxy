package utils

import (
	"net/url"
	"regexp"
)

// FormatURL - find url in query param
func FormatURL(u string) *url.URL {
	s := regexp.MustCompile(`(?:/query\?q=)(.{0,})`)
	r := string(s.ReplaceAll([]byte(u), []byte("$1")))
	formattedURL, err := url.Parse(r)
	if err != nil {
		panic(err)
	}
	return formattedURL
}
