package domain

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/joshbatley/proxy/utils"
)

// Record Request Data struct
type Record struct {
	URL        string
	Body       []byte
	Headers    string
	Status     int
	Method     string
	Datetime   time.Time
	Collection int64
}

// NewRecord take raw formats them read for sql saving
func NewRecord(u *url.URL, b io.ReadCloser, h http.Header, s int, m string, c int64) *Record {
	buf := new(bytes.Buffer)
	buf.ReadFrom(b)

	return &Record{
		URL:        u.String(),
		Body:       buf.Bytes(),
		Headers:    utils.HeadersToString(h),
		Status:     s,
		Method:     m,
		Datetime:   time.Now(),
		Collection: c,
	}

}
