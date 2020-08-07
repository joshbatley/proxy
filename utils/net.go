package utils

import "github.com/sparrc/go-ping"

// Connected if connected to the internet
func Connected() (ok bool) {
	_, err := ping.NewPinger("www.google.com")
	if err != nil {
		return true
	}
	return
}
