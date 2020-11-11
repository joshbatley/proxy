package connection

import (
	"net"
)

// IsOnline calls google to check for a connection
func IsOnline(urls []string) bool {
	if len(urls) == 0 {
		urls = append(urls, "google.com")
	}

	for _, u := range urls {
		ips, _ := net.LookupIP(u)
		if len(ips) > 0 {
			return true
		}
	}
	return false
}
