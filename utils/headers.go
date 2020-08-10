package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

// Cors applies cors headers
func Cors(h http.Header) {
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Access-Control-Allow-Methods", "*")
	h.Set("Access-Control-Allow-Headers", "*")
}

// HeadersToString takes all the headers and maps them to a string
func HeadersToString(h http.Header) string {
	b := new(bytes.Buffer)
	for k, v := range h {
		fmt.Fprintf(b, "%s|%s\n", k, strings.Join(v, " "))
	}

	return b.String()
}

// StringToHeaders take HeadersToString and reverts its
func StringToHeaders(h string, w http.ResponseWriter) {
	for _, i := range strings.Split(h, "\n") {
		h := strings.Split(i, "|")
		if len(h) >= 2 {
			w.Header().Set(h[0], h[1])
		}
	}
}
