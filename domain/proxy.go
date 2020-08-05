package domain

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Proxy struct {
	ModifyResponse func(*http.Response) error
	URL            *url.URL
}

func (p *Proxy) Serve(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, CONNECT, OPTIONS, TRACE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "*, Origin, X-Requested-With, Content-Type, Accept")
	}

	reverseProxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = p.URL.Scheme
			req.URL.Host = p.URL.Host
			req.URL.Path = p.URL.Path
			req.Host = p.URL.Host
			req.URL.RawQuery = p.URL.RawQuery
		},
		ModifyResponse: p.ModifyResponse,
	}

	log.Println("Method:", r.Method, "Calling:", p.URL.Host)

	reverseProxy.ServeHTTP(w, r)
}
