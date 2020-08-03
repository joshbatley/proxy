package handler

import (
	"log"
	"net/http"

	"github.com/joshbatley/proxy/def"
	"github.com/joshbatley/proxy/service"
)

var cache []def.Record

// QueryServe - as
func QueryServe(w http.ResponseWriter, r *http.Request) {
	query := def.NewQuery(w, r, service.ModifyResponse)

	if re := service.GetPreResponse(query.URL, r, w); re {
		return
	}

	if re := service.SendCache(query.URL, w, cache); re {
		return
	}

	log.Println("no entry fetching from server")

	query.ReverseProxy.ServeHTTP(w, r)
}
