package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joshbatley/proxy/domain/collections"
	"github.com/joshbatley/proxy/domain/endpoints"
	"github.com/joshbatley/proxy/domain/responses"
	"github.com/joshbatley/proxy/domain/rules"
	"github.com/joshbatley/proxy/internal/connection"
	"github.com/joshbatley/proxy/internal/engine"
	"github.com/joshbatley/proxy/internal/fail"
	"github.com/joshbatley/proxy/internal/params"
	"go.uber.org/zap"
)

// QueryHandler Http handler for any query response
type QueryHandler struct {
	collections *collections.Manager
	endpoints   *endpoints.Manager
	responses   *responses.Manager
	rules       *rules.Manager
	engine      *engine.RuleEngine
	log         *zap.SugaredLogger
	d           string
}

// NewQueryHandler constructs a new QueryHandler
func NewQueryHandler(
	collections *collections.Manager,
	endpoints *endpoints.Manager,
	responses *responses.Manager,
	rules *rules.Manager,
	log *zap.SugaredLogger,
) QueryHandler {
	return QueryHandler{
		collections: collections,
		endpoints:   endpoints,
		responses:   responses,
		rules:       rules,
		engine:      &engine.RuleEngine{},
		log:         log,
	}
}

// Serve Sets up all the logic for a reverse proxy and save and sends cached versions
func (q QueryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params, err := params.Parse(mux.Vars(r), r.URL)
	if err != nil {
		q.log.Error("Param parse fail")
		badRequest(err, w)
		return
	}

	// Check the collection exist (if not default)
	_, err = q.collections.Get(params.Collection)
	if err == fail.ErrNoData {
		q.log.Warn("No collection found")
		badRequest(fail.MissingColErr(err), w)
		return
	}

	// With colleciton load rules for store
	rules, err := q.rules.Get(params.Collection)
	if err != nil {
		q.log.Warn("fail to get rules")
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

	// All to save
	if !q.engine.CheckStore() {
		q.log.Info("Proxing", params.QueryURL)
		q.reverseProxy(w, r, params, func(re *http.Response) error {
			if q.engine.EnableCors() {
				corsHeaders(re.Header)
			}
			return nil
		})
	}

	// Check for endpoint
	endpointID, err := q.endpoints.GetOrCreate(params.QueryURL.String(), r.Method, params.Collection)

	if err != nil {
		badRequest(err, w)
		return
	}

	// return cache
	found, err := q.returnCache(w, r, params, endpointID)
	if found {
		return
	}
	if err != nil {
		badRequest(err, w)
		return
	}

	q.reverseProxy(w, r, params, func(re *http.Response) error {
		// Depulicate the body to reapply to response later
		buf, _ := ioutil.ReadAll(re.Body)
		q.log.Info("saving response for", re.Request.URL)

		err := q.saveResponse(
			q.d,
			re.Request.URL,
			ioutil.NopCloser(bytes.NewBuffer(buf)),
			re.Header,
			re.StatusCode,
			re.Request.Method,
			endpointID,
		)

		if err != nil {
			q.log.Error("Failed to save response")

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
	})

}

func (q *QueryHandler) returnCache(
	w http.ResponseWriter, r *http.Request, p *params.Params, e uuid.UUID,
) (bool, error) {

	res, err := q.responses.Get(
		p.QueryURL.String(),
		e.String(),
		r.Method,
	)

	if err == fail.ErrNoData {
		q.log.Info("no data found proxy request")
		return false, nil
	}

	if err != nil {
		q.log.Error("Getting response failed")
		return false, err
	}

	q.d = res.ID
	if q.engine.HasExpired(res.DateTime) {
		q.log.Info("response has expired - refresh data")
		return false, nil
	}

	q.log.Info("returned saved response")
	// Headers from string to headers
	for _, i := range strings.Split(res.Headers, "\n") {
		h := strings.Split(i, "|")
		if len(h) >= 2 {
			w.Header().Set(h[0], h[1])
		}
	}

	w.Header().Set("x-Proxy", "served from cache")
	w.WriteHeader(res.Status)
	w.Write(res.Body)
	return true, nil

}

func (q *QueryHandler) reverseProxy(
	w http.ResponseWriter, r *http.Request, p *params.Params, mr func(r *http.Response) error,
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

	reverseProxy := httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: mr,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			if !connection.IsOnline(nil) {
				badRequest(fail.OfflineError(err), w)
			} else {
				q.log.Warn("Internal Error on reverse Proxy - ", err)
				badRequest(fail.InternalError(err), w)
			}
		},
	}

	reverseProxy.ServeHTTP(w, r)
}

func (q *QueryHandler) saveResponse(id string, u *url.URL, b io.ReadCloser, h http.Header, s int, m string, e uuid.UUID) error {
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
