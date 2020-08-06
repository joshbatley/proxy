package handler

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/joshbatley/proxy/domain"
	"github.com/joshbatley/proxy/repository"
	"github.com/joshbatley/proxy/utils"
)

// QueryHandler -
type QueryHandler struct {
	CacheRepository *repository.CacheRepository
}

// Serve -
func (q *QueryHandler) Serve(w http.ResponseWriter, r *http.Request) {
	url := utils.FormatURL(r.URL.String())

	d, err := q.CacheRepository.GetCache(url.String())

	if err == nil {
		log.Println("served from cache")
		q.sendCache(d, w)
		return
	}

	p := domain.Proxy{
		ModifyResponse: q.saveReponse,
		URL:            url,
	}

	p.Serve(w, r)
}

// SaveReponse -
func (q *QueryHandler) saveReponse(r *http.Response) error {
	r.Header.Set("Access-Control-Allow-Origin", "*")
	r.Header.Set("Access-Control-Allow-Methods", "*")
	r.Header.Set("Access-Control-Allow-Headers", "*")

	// Depulicate the body to reapply to response later
	buf, _ := ioutil.ReadAll(r.Body)
	newC := domain.NewRecord(
		r.Request.URL,
		ioutil.NopCloser(bytes.NewBuffer(buf)),
		r.Header,
		r.StatusCode,
		r.Request.Method,
	)

	err := q.CacheRepository.SaveCache(newC)
	if err != nil {
		log.Println(err)
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

	return nil
}

// SendCache -
func (q *QueryHandler) sendCache(d repository.Cache, w http.ResponseWriter) {
	for _, i := range strings.Split(d.Headers, "\n") {
		h := strings.Split(i, "|")
		if len(h) >= 2 {
			k := h[0]
			v := h[1]
			w.Header().Set(k, v)
		}
	}
	w.Header().Set("x-Proxy", "served from cache")
	w.WriteHeader(d.Status)
	w.Write(d.Body)
}
