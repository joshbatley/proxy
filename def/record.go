package def

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

// Record - Request Data struct
type Record struct {
	URL     *url.URL
	Body    []byte
	Headers http.Header
	Status  int
	Method  string
}

// URLString - Returns url as string
func (r *Record) URLString() string {
	return r.URL.String()
}

// HeadersToString - Returns Headers as 'foo=[bar,baz];'
func (r *Record) HeadersToString() string {
	b := new(bytes.Buffer)
	for key, value := range r.Headers {
		fmt.Fprintf(b, "%s=%s;\n", key, value)
	}

	return b.String()
}

// StringToHeader - Takes string and formats to Headers
func (r *Record) StringToHeader(str string) {
	// Split by semicolon, values by array
	//r.Headers.Set(key string, value string)
}
