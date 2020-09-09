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

// ResponseRow returns struct from the database
type ResponseRow struct {
	ID     int    `db:"ID"`
	Status int    `db:"Status"`
	URL    string `db:"URL"`
	// Returns Headers as 'foo=bar; baz, other \n'
	Headers  string `db:"Headers"`
	Body     []byte `db:"Body"`
	DateTime int64  `db:"DateTime"`
}

// Collection returns struct from the database
type Collection struct {
	ID   int64  `db:"ID"`
	Name string `db:"Name"`
}

// Rule returns a single rule
type Rule struct {
	Pattern      string `db:"Pattern"`
	SaveResponse int    `db:"SaveResponse"`
	ForceCors    int    `db:"ForceCors"`
	Expiry       int    `db:"Expiry"`
}

// Endpoint returns a single endpoint
type Endpoint struct {
	ID     int64  `db:"ID"`
	Status int    `db:"PreferedStatus"`
	Method string `db:"Method"`
	URL    string `db:"URL"`
}

// Response request data struct
type Response struct {
	URL      string
	Body     []byte
	Headers  string
	Status   int
	Method   string
	DateTime int64
	Endpoint int64
}

// NewResponse take raw formats them parses them for sql saving
func NewResponse(u *url.URL, b io.ReadCloser, h http.Header, s int, m string, e int64) *Response {
	buf := new(bytes.Buffer)
	buf.ReadFrom(b)

	headers := new(bytes.Buffer)
	for k, v := range h {
		fmt.Fprintf(headers, "%s|%s\n", k, strings.Join(v, " "))
	}

	return &Response{
		URL:      u.String(),
		Body:     buf.Bytes(),
		Headers:  headers.String(),
		Status:   s,
		Method:   m,
		DateTime: time.Now().Unix(),
		Endpoint: e,
	}
}
