package handler

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/joshbatley/proxy/def"
	"github.com/joshbatley/proxy/service"
)

var cache []def.Record

func handleOptions(w http.ResponseWriter, m string) bool {
	if m == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, CONNECT, OPTIONS, TRACE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "*, Origin, X-Requested-With, Content-Type, Accept")
		w.WriteHeader(http.StatusNoContent)
		return true
	}
	return false
}

func director(url *url.URL) func(req *http.Request) {
	return func(req *http.Request) {
		req.URL.Scheme = url.Scheme
		req.URL.Host = url.Host
		req.URL.Path = url.Path
		req.Host = url.Host
		req.URL.RawQuery = url.RawQuery
	}
}

// QueryServe - as
func QueryServe(w http.ResponseWriter, r *http.Request) {
	url := service.FormatURL(r.URL.String())

	if re := handleOptions(w, r.Method); re {
		return
	}

	if re := service.GetPreResponse(url, r, w); re {
		return
	}

	if re := service.SendCache(url, w, cache); re {
		return
	}

	b := new(bytes.Buffer)
	for key, value := range r.Header {
		fmt.Fprintf(b, "%s=%s,\n", key, value)
	}
	log.Println(b.String())

	log.Println("no entry fetching from server")

	reverseProxy := httputil.ReverseProxy{
		Director:       director(url),
		ModifyResponse: service.ModifyResponse,
	}

	reverseProxy.ServeHTTP(w, r)
}
