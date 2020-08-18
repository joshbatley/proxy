package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joshbatley/proxy/domain"
	"github.com/joshbatley/proxy/repository"
	"github.com/joshbatley/proxy/utils"
)

// QueryHandler Http handler for any query response
type QueryHandler struct {
	CacheRepository *repository.CacheRepository
}

var collection int64

// Serve Sets up all the logic for a reverse proxy and save and sends cached versions
func (q *QueryHandler) Serve(w http.ResponseWriter, r *http.Request) {
	params, err := utils.ParseParams(mux.Vars(r), r.URL)
	collection = params.Collection

	if err != nil {
		w.Header().Set("Content-Type", "application/json, text/plain, */*")
		w.WriteHeader(http.StatusBadRequest)
		jsonString, _ := json.Marshal(err)
		w.Write([]byte(jsonString))
		return
	}

	d, err := q.CacheRepository.GetCache(params.QueryURL.String(), collection)
	switch {
	case err == nil:
		log.Println("served from cache")
		q.sendCache(d, w)
		return
	case err != sql.ErrNoRows:
		log.Fatal(err)
	}

	p := domain.Proxy{
		ModifyResponse: q.saveReponse,
		URL:            params.QueryURL,
	}

	p.Proxy(w, r)
}

func (q *QueryHandler) saveReponse(r *http.Response) error {
	// Apply headers to skip inbuild security
	utils.Cors(r.Header)
	// Depulicate the body to reapply to response later
	buf, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	err := q.CacheRepository.SaveCache(
		domain.NewRecord(
			r.Request.URL,
			ioutil.NopCloser(bytes.NewBuffer(buf)),
			r.Header,
			r.StatusCode,
			r.Request.Method,
			collection,
		),
	)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (q *QueryHandler) sendCache(d *repository.CacheRow, w http.ResponseWriter) {
	utils.StringToHeaders(d.Headers, w)
	w.Header().Set("x-Proxy", "served from cache")
	w.WriteHeader(d.Status)
	w.Write(d.Body)
}
