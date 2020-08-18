package utils

import "github.com/sparrc/go-ping"

// Connected check if the outside internet is reachable
func Connected() (ok bool) {
	_, err := ping.NewPinger("www.google.com")
	if err != nil {
		return true
	}
	return
}
