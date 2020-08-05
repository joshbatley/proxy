package handler

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/joshbatley/proxy/domain"
	"github.com/joshbatley/proxy/repository"
)

// QueryHandler -
type QueryHandler struct {
	CacheRepository *repository.CacheRepository
}

// Serve -
func (q *QueryHandler) Serve(w http.ResponseWriter, r *http.Request) {
	url := formatURL(r.URL.String())

	d, err := q.CacheRepository.GetCache(url.String())

	if err == nil {
		q.sendCache(d, w)
		return
	}

	p := domain.Proxy{
		ModifyResponse: q.SaveReponse,
		URL:            url,
	}

	p.Serve(w, r)
}

// SaveReponse -
func (q *QueryHandler) SaveReponse(res *http.Response) error {
	// Depulicate the body to reapply to response later
	buf, _ := ioutil.ReadAll(res.Body)

	newC := domain.NewRecord(
		res.Request.URL,
		ioutil.NopCloser(bytes.NewBuffer(buf)),
		res.Header,
		res.StatusCode,
		res.Request.Method,
	)

	q.CacheRepository.SaveCache(newC)
	res.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

	return nil
}

// SendCache -
func (q *QueryHandler) sendCache(d repository.Cache, w http.ResponseWriter) {
	log.Println("found in cache sending cache")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(d.Status)
	w.Write(d.Body)
}

func formatURL(u string) *url.URL {
	s := regexp.MustCompile(`(?:/query\?q=)(.{0,})`)
	r := string(s.ReplaceAll([]byte(u), []byte("$1")))

	formattedURL, err := url.Parse(r)

	if err != nil {
		panic(err)
	}

	return formattedURL
}
