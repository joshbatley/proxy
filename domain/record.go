package domain

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Record - Request Data struct
type Record struct {
	URL        string
	Body       []byte
	Headers    string
	Status     int
	Method     string
	Datetime   time.Time
	Collection string
}

// URLString - Returns url as strin
// func (r *Record) URLString() string {
// 	return r.URL.String()
// }

// HeadersToString - Returns Headers as 'foo=[bar,baz];'
func HeadersToString(h http.Header) string {
	b := new(bytes.Buffer)
	for k, v := range h {
		fmt.Fprintf(b, "%s|%s\n", k, strings.Join(v, " "))
	}

	return b.String()
}

// StringToHeader - Takes string and formats to Headers
func (r *Record) StringToHeader(str string) {
	// Split by semicolon, values by array
	//r.Headers.Set(key string, value string)
}

// NewRecord -
func NewRecord(u *url.URL, b io.ReadCloser, h http.Header, s int, m string) Record {
	buf := new(bytes.Buffer)
	buf.ReadFrom(b)

	return Record{
		URL:        u.String(),
		Body:       buf.Bytes(),
		Headers:    HeadersToString(h),
		Status:     s,
		Method:     m,
		Datetime:   time.Now(),
		Collection: "",
	}

}
