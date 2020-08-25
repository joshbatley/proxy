package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Record request data struct
type Record struct {
	URL        string
	Body       []byte
	Headers    string
	Status     int
	Method     string
	Datetime   time.Time
	Collection int64
}

// CacheRow returns struct from the database
type CacheRow struct {
	ID     int
	Status int
	URL    string
	// Returns Headers as 'foo=bar; baz, other \n'
	Headers string
	Body    []byte
}

func headersToString(h http.Header) string {
	b := new(bytes.Buffer)
	for k, v := range h {
		fmt.Fprintf(b, "%s|%s\n", k, strings.Join(v, " "))
	}

	return b.String()
}

// NewRecord take raw formats them parses them for sql saving
func NewRecord(u *url.URL, b io.ReadCloser, h http.Header, s int, m string, c int64) *Record {
	buf := new(bytes.Buffer)
	buf.ReadFrom(b)

	return &Record{
		URL:        u.String(),
		Body:       buf.Bytes(),
		Headers:    headersToString(h),
		Status:     s,
		Method:     m,
		Datetime:   time.Now(),
		Collection: c,
	}
}
