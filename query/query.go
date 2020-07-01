package query

import (
	"fmt"
	"goproxy/utils"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var cache []Req

func director(url *url.URL) func(req *http.Request) {
	return func(req *http.Request) {
		req.URL.Scheme = url.Scheme
		req.URL.Host = url.Host
		req.URL.Path = url.Path
		req.Host = url.Host
		req.URL.RawQuery = url.RawQuery
	}
}

// Serve -
func Serve(w http.ResponseWriter, r *http.Request) {
	url := utils.FormatURL(r.URL.String())

	if re := handleOptions(w, r.Method); re {
		return
	}

	if re := getPreResponse(url, r, w); re {
		return
	}

	if re := sendCache(url, w, cache); re {
		return
	}

	fmt.Println("no entry fetching from server")

	reverseProxy := httputil.ReverseProxy{
		Director:       director(url),
		ModifyResponse: modifyResponse(cache),
	}

	reverseProxy.ServeHTTP(w, r)
}
