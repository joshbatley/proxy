package def

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
)

// Query -
type Query struct {
	URL          *url.URL
	ReverseProxy *httputil.ReverseProxy
}

func formatURL(u string) *url.URL {
	s := regexp.MustCompile(`(?:/query\?q=)(.{0,})`)
	r := string(s.ReplaceAll([]byte(u), []byte("$1")))

	formattedURL, err := url.Parse(r)

	if err != nil {
		panic(err)
	}

	return formattedURL
}

// NewQuery -
func NewQuery(w http.ResponseWriter, r *http.Request, mr func(*http.Response) error) Query {
	url := formatURL(r.URL.String())

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, CONNECT, OPTIONS, TRACE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "*, Origin, X-Requested-With, Content-Type, Accept")
	}

	reverseProxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = url.Scheme
			req.URL.Host = url.Host
			req.URL.Path = url.Path
			req.Host = url.Host
			req.URL.RawQuery = url.RawQuery
		},
		ModifyResponse: mr,
	}

	return Query{
		URL:          url,
		ReverseProxy: &reverseProxy,
	}
}
