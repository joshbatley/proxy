package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joshbatley/proxy"
	"github.com/joshbatley/proxy/internal/engine"
	"github.com/joshbatley/proxy/internal/store"
	"github.com/joshbatley/proxy/internal/utils"
)

// QueryHandler Http handler for any query response
type QueryHandler struct {
	Store      *store.Store
	collection int64
}

// Serve Sets up all the logic for a reverse proxy and save and sends cached versions
func (q QueryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params, err := utils.ParseParams(mux.Vars(r), r.URL)
	if err != nil {
		badRequest(err, w)
		return
	}
	q.collection = params.Collection
	if _, err := q.Store.GetCollection(q.collection); err == sql.ErrNoRows {
		badRequest(proxy.MissingColErr(err), w)
		return
	}

	res, err := q.Store.GetRules(q.collection)
	if err != nil {
		badRequest(proxy.MissingColErr(err), w)
	}

	state := engine.Engine(res, params.QueryURL)
	switch state {
	case engine.StateSaving:
		d, err := q.Store.GetCache(params.QueryURL.String(), q.collection)
		if err != nil {
			log.Fatal("DB Fell over", err)
		}
		if d != nil {
			log.Println("served from cache")
			sendCache(d, w)
			return
		}
	}

	log.Printf("Getting new data and will save %v ", state)
	q.reverseProxy(params.QueryURL, w, r, state)

}

func (q *QueryHandler) saveReponse(r *http.Response) error {
	log.Print("Saving response")
	// Apply headers to skip inbuild security
	corsHeaders(r.Header)

	// Depulicate the body to reapply to response later
	buf, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	if err := q.Store.SaveCache(
		proxy.NewRecord(
			r.Request.URL,
			ioutil.NopCloser(bytes.NewBuffer(buf)),
			r.Header,
			r.StatusCode,
			r.Request.Method,
			q.collection,
		),
	); err != nil {
		log.Println(err)
	}

	return nil
}

func (q *QueryHandler) reverseProxy(URL *url.URL, w http.ResponseWriter, r *http.Request, state engine.State) {
	// Always allows cors, all webapps to bypass security
	if r.Method == http.MethodOptions {
		corsHeaders(w.Header())
		return
	}

	reverseProxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Header.Del("Origin")
			req.Header.Del("Referer")
			req.URL.Scheme = URL.Scheme
			req.URL.Host = URL.Host
			req.URL.Path = URL.Path
			req.Host = URL.Host
			req.URL.RawQuery = URL.RawQuery
		},
		ModifyResponse: func(r *http.Response) error {
			if state == engine.StateSaving {
				log.Println(state)
				return q.saveReponse(r)
			}
			return nil
		},
	}

	log.Println("Method:", r.Method, "Calling:", URL.String())

	reverseProxy.ServeHTTP(w, r)
}

func sendCache(d *proxy.CacheRow, w http.ResponseWriter) {
	for _, i := range strings.Split(d.Headers, "\n") {
		h := strings.Split(i, "|")
		if len(h) >= 2 {
			w.Header().Set(h[0], h[1])
		}
	}

	w.Header().Set("x-Proxy", "served from cache")
	w.WriteHeader(d.Status)
	w.Write(d.Body)
}

func badRequest(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json, text/plain, */*")
	w.WriteHeader(http.StatusBadRequest)
	jsonString, _ := json.Marshal(err)
	w.Write(jsonString)
}

func corsHeaders(h http.Header) {
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Access-Control-Allow-Methods", "*")
	h.Set("Access-Control-Allow-Headers", "*")
}
