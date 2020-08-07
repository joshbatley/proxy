package handler

import (
	"bytes"
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"

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
	} else if err != sql.ErrNoRows {
		log.Fatal(err)
	}

	p := domain.Proxy{
		ModifyResponse: q.saveReponse,
		URL:            url,
	}

	p.Serve(w, r)
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
