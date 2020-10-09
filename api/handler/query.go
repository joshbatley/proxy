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
	d           string
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

	// Check the collection exist (if not default)
	_, err = q.collections.Get(params.Collection)
	if err == sql.ErrNoRows {
		log.Println("No collection found")
		badRequest(fail.MissingColErr(err), w)
		return
	}

	// With colleciton load rules for store
	rules, err := q.rules.Get(params.Collection)
	if err != nil {
		log.Println("fail to get rules")
		badRequest(err, w)
		return
	}

	// Convert rules type
	engineRules := make([]engine.Rule, len(rules))
	for _, v := range rules {
		engineRules = append(engineRules, engine.Rule(v))
	}

	// pass rules to engine
	q.engine.LoadRules(params.QueryURL, params.Collection, engineRules)

	// Check if method is OPTIONS and if Engine need to override
	if r.Method == http.MethodOptions && q.engine.EnableCors() {
		corsHeaders(w.Header())
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var endpointID string

	// All to save
	if q.engine.CheckStore() {
		// Check for endpoint
		end, err := q.endpoints.Get(params.QueryURL.String(), r.Method, params.Collection)

		if err == sql.ErrNoRows {
			log.Println("Created enpoint")
			id, err := q.endpoints.Save(params.Collection, params.QueryURL.String(), r.Method)
			endpointID = id.String()
			log.Println(id, err)
		}

		if err != nil && err != sql.ErrNoRows {
			log.Println("Endpoints didnt save")
			badRequest(err, w)
			return
		}
		if err == nil {
			endpointID = end.ID
		}

		// return cache
		ok, err := q.returnCache(w, r, params, endpointID)
		if ok {
			return
		}
		if err != nil {
			badRequest(err, w)
			return
		}
	}

	// Allows fallthrough and
	log.Println("Proxing", params.QueryURL)
	q.reverseProxy(w, r, params, endpointID)
}

func (q *QueryHandler) returnCache(
	w http.ResponseWriter, r *http.Request, p *params.Params, e string,
) (bool, error) {

	d, err := q.responses.Get(
		p.QueryURL.String(),
		e,
		r.Method,
	)

	if err == sql.ErrNoRows {
		log.Println("no data found proxy request")
		return false, nil
	}

	if err != nil && err != sql.ErrNoRows {
		log.Println("Getting response failed")
		return false, err
	}

	q.d = d.ID
	if q.engine.HasExpired(d.DateTime) {
		log.Println("response has expired - refresh data")
		return false, nil
	}

	log.Println("returned saved response")
	// Headers from string to headers
	for _, i := range strings.Split(d.Headers, "\n") {
		h := strings.Split(i, "|")
		if len(h) >= 2 {
			w.Header().Set(h[0], h[1])
		}
	}

	w.Header().Set("x-Proxy", "served from cache")
	w.WriteHeader(d.Status)
	w.Write(d.Body)
	return true, nil

}

func (q *QueryHandler) reverseProxy(
	w http.ResponseWriter, r *http.Request, p *params.Params, e string,
) {
	director := func(req *http.Request) {
		req.Header.Del("Origin")
		req.Header.Del("Referer")
		req.URL.Scheme = p.QueryURL.Scheme
		req.URL.Host = p.QueryURL.Host
		req.URL.Path = p.QueryURL.Path
		req.Host = p.QueryURL.Host
		req.URL.RawQuery = p.QueryURL.RawQuery
	}

	modifyResponse := func(re *http.Response) error {
		if !q.engine.CheckStore() {
			return nil
		}

		// Depulicate the body to reapply to response later
		buf, _ := ioutil.ReadAll(re.Body)

		err := q.saveResponse(
			q.d,
			re.Request.URL,
			ioutil.NopCloser(bytes.NewBuffer(buf)),
			re.Header,
			re.StatusCode,
			re.Request.Method,
			e,
		)

		if err != nil {
			log.Println("Failed to save response")

			re.Header = http.Header{}
			re.Header.Set("Content-Type", "application/json, text/plain, */*")
			re.StatusCode = http.StatusBadRequest
			jsonString, _ := json.Marshal(fail.InternalError(err))
			re.Body = ioutil.NopCloser(bytes.NewBuffer(jsonString))
			return nil
		}

		// Apply headers to skip inbuild security
		if q.engine.EnableCors() {
			corsHeaders(re.Header)
		}

		re.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
		return nil
	}

	reverseProxy := httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: modifyResponse,
	}

	reverseProxy.ServeHTTP(w, r)
}

func (q *QueryHandler) saveResponse(id string, u *url.URL, b io.ReadCloser, h http.Header, s int, m string, e string) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(b)

	headers := new(bytes.Buffer)
	for k, v := range h {
		fmt.Fprintf(headers, "%s|%s\n", k, strings.Join(v, " "))
	}

	return q.responses.Save(
		id,
		u.String(),
		headers.String(),
		buf.Bytes(),
		s,
		m,
		e,
	)
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
