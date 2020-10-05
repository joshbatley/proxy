package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joshbatley/proxy/domain/collections"
	"github.com/joshbatley/proxy/domain/endpoints"
	"github.com/joshbatley/proxy/domain/responses"
	"github.com/joshbatley/proxy/domain/rules"
	"github.com/joshbatley/proxy/internal/engine"
	"github.com/joshbatley/proxy/internal/fail"
	"github.com/joshbatley/proxy/internal/params"
)

// QueryHandler Http handler for any query response
type QueryHandler struct {
	collections *collections.Manager
	endpoints   *endpoints.Manager
	responses   *responses.Manager
	rules       *rules.Manager
	engine      *engine.RuleEngine
}

// NewQueryHandler constructs a new QueryHandler
func NewQueryHandler(
	collections *collections.Manager,
	endpoints *endpoints.Manager,
	responses *responses.Manager,
	rules *rules.Manager,
) QueryHandler {
	return QueryHandler{
		collections: collections,
		endpoints:   endpoints,
		responses:   responses,
		rules:       rules,
		engine:      &engine.RuleEngine{},
	}
}

// Serve Sets up all the logic for a reverse proxy and save and sends cached versions
func (q QueryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params, err := params.Parse(mux.Vars(r), r.URL)
	if err != nil {
		log.Println("Param parse fail")
		badRequest(err, w)
		return
	}

	_, err = q.collections.Get(params.Collection)
	if err == sql.ErrNoRows {
		log.Println("No collection found")
		badRequest(fail.MissingColErr(err), w)
		return
	}

	rules, err := q.rules.Get(params.Collection)
	if err != nil {
		log.Println("fail to get rules")
		badRequest(err, w)
		return
	}

	q.engine.LoadRules(params, rules)

	if r.Method == http.MethodOptions && q.engine.EnableCors() {
		corsHeaders(w.Header())
		w.WriteHeader(http.StatusNoContent)
		return
	}

	state := q.engine.GetState()

	e, err := q.endpoints.GetOrSave(params.QueryURL.String(), r.Method, params.Collection)
	if err != nil {
		log.Println("Endpoints didnt save/find")
		badRequest(err, w)
		return
	}

	if state == engine.StateSaving {
		d, err := q.responses.Get(
			params.QueryURL.String(),
			e.ID,
			r.Method,
		)
		if err != nil {
			log.Println("Getting response failed")
			badRequest(err, w)
			return
		}
		if d != nil && !q.engine.HasExpired(d.DateTime) {
			sendResponse(d, w)
			return
		}
	}

	q.reverseProxy(w, r, params)
}

func (q QueryHandler) reverseProxy(w http.ResponseWriter, r *http.Request, p *params.Params) {
	d := func(req *http.Request) {
		req.Header.Del("Origin")
		req.Header.Del("Referer")
		req.URL.Scheme = p.QueryURL.Scheme
		req.URL.Host = p.QueryURL.Host
		req.URL.Path = p.QueryURL.Path
		req.Host = p.QueryURL.Host
		req.URL.RawQuery = p.QueryURL.RawQuery
	}

	mr := func(re *http.Response) error {
		if state := q.engine.GetState(); state == engine.StateSaving {
			// Depulicate the body to reapply to response later
			buf, _ := ioutil.ReadAll(re.Body)
			err := q.saveResponse(re.Request.URL,
				ioutil.NopCloser(bytes.NewBuffer(buf)),
				re.Header,
				re.StatusCode,
				re.Request.Method,
				p.Collection,
			)

			if err != nil {
				log.Println("Failed to save response")
				badResponse(fail.InternalError(err), re)
				return nil
			}
			// Apply headers to skip inbuild security
			if q.engine.EnableCors() {
				corsHeaders(re.Header)
			}
			re.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
			return nil
		}
		return nil
	}

	reverseProxy := httputil.ReverseProxy{
		Director:       d,
		ModifyResponse: mr,
	}

	reverseProxy.ServeHTTP(w, r)
}

func (q QueryHandler) saveResponse(u *url.URL, b io.ReadCloser, h http.Header, s int, m string, e int64) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(b)

	headers := new(bytes.Buffer)
	for k, v := range h {
		fmt.Fprintf(headers, "%s|%s\n", k, strings.Join(v, " "))
	}

	return q.responses.Save(
		u.String(),
		headers.String(),
		buf.Bytes(),
		s,
		m,
		e,
	)
}

func sendResponse(d *responses.Response, w http.ResponseWriter) {
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

func badResponse(err error, r *http.Response) {
	r.Header = http.Header{}
	r.Header.Set("Content-Type", "application/json, text/plain, */*")
	r.StatusCode = http.StatusBadRequest
	jsonString, _ := json.Marshal(err)

	if len(jsonString) == 2 {
		jsonString, _ = json.Marshal(fail.InternalError(err))
	}

	r.Body = ioutil.NopCloser(bytes.NewBuffer(jsonString))
}

func badRequest(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json, text/plain, */*")
	w.WriteHeader(http.StatusBadRequest)
	jsonString, _ := json.Marshal(err)

	if len(jsonString) == 2 {
		jsonString, _ = json.Marshal(fail.InternalError(err))
	}

	w.Write(jsonString)
}

func corsHeaders(h http.Header) {
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Access-Control-Allow-Methods", "*")
	h.Set("Access-Control-Allow-Headers", "*")
}
