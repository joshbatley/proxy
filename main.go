package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"goproxy/response"
	"goproxy/utils"
)

var cache []response.Req

func query(w http.ResponseWriter, r *http.Request) {
	url := utils.FormatURL(r.URL.String())

	if re := response.HandleOptions(w, r.Method); re {
		return
	}

	if re := response.GetPreResponse(url, r, w); re {
		return
	}

	if re := response.SendCache(url, w, cache); re {
		return
	}

	fmt.Println("no entry fetching from server")

	reverseProxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = url.Scheme
			req.URL.Host = url.Host
			req.URL.Path = url.Path
			req.Host = url.Host
			req.URL.RawQuery = url.RawQuery
		},
		ModifyResponse: func(res *http.Response) error {
			return response.ModifyResponse(res, cache)
		},
	}

	reverseProxy.ServeHTTP(w, r)
}

func main() {
	config, err := utils.ReadConfig("./config.yml")

	if err != nil {
		panic("Config unreadable")
	}

	http.HandleFunc("/query", query)
	http.ListenAndServe(":"+config.Port, nil)

	fmt.Println("listing on 127.0.0.1:" + config.Port)
}
