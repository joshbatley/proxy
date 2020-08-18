package domain

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/joshbatley/proxy/utils"
)

// Proxy required info to set up a proxy request
type Proxy struct {
	ModifyResponse func(*http.Response) error
	URL            *url.URL
}

// Proxy setups and return the reverse proxy request
func (p *Proxy) Proxy(w http.ResponseWriter, r *http.Request) {
	// Always allows cors, all webapps to bypass security
	if r.Method == http.MethodOptions {
		utils.Cors(w.Header())
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
