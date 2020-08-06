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

func preflight(w http.ResponseWriter, r *http.Request) (ok bool) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, CONNECT, OPTIONS, TRACE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "*, Origin, X-Requested-With, Content-Type, Accept")
		return true
	}
	return
}

func (p *Proxy) Serve(w http.ResponseWriter, r *http.Request) {
	if ok := preflight(w, r); ok {
		return
	}

	reverseProxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Header.Del("Origin")
			req.Header.Del("Referer")
			req.URL.Scheme = p.URL.Scheme
			req.URL.Host = p.URL.Host
			req.URL.Path = p.URL.Path
			req.Host = p.URL.Host
			req.URL.RawQuery = p.URL.RawQuery
		},
		ModifyResponse: p.ModifyResponse,
	}

	log.Println("Method:", r.Method, "Calling:", p.URL.String())

	reverseProxy.ServeHTTP(w, r)
}
