package connection

import (
	"net/http"
)

// IsOffline calls google to check for a connection
func IsOffline(urls []string) bool {
	if len(urls) == 0 {
		urls = append(urls, "http://clients3.google.com/generate_204")
	}

	for _, u := range urls {
		_, err := http.Get(u)
		if err == nil {
			return false
		}
	}
	return false
}
