package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
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
	log         *zap.SugaredLogger
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
		log:         log,
	}
}

type response struct {
	headers string
	status  int
	body    []byte
}

type ids struct {
	endpoint uuid.UUID
	id       uuid.UUID
}

// Serve Sets up all the logic for a reverse proxy and save and sends cached versions
func (q QueryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params, err := params.Parse(mux.Vars(r), r.URL)
	if err != nil {
		q.log.Error("Param parse failed")
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

	engine, err := q.loadEngine(params)
	if err != nil {
		q.log.Warn("Failed to load rules")
		badRequest(err, w)
		return
	}

	// Check if method is OPTIONS and if Engine need to override
	if r.Method == http.MethodOptions && engine.EnableCors() {
		q.log.Info("Cors request")
		corsHeaders(w.Header())
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Skip Store check, and straight proxy
	if !engine.CheckStore() {
		q.log.Info("Skipping cache check and proxing: ", params.QueryURL)
		reverseProxy(w, r, params, func(re *http.Response) error {
			if engine.EnableCors() {
				corsHeaders(re.Header)
			}
			return nil
		}, q.log)
		return
	}

	cachedChannel := make(chan response)
	proxyChannel := make(chan ids)
	errChannel := make(chan error)

	defer close(cachedChannel)
	defer close(proxyChannel)
	defer close(errChannel)

	go q.checkResponses(params, r, engine, cachedChannel, proxyChannel, errChannel)

	select {
	case err := <-errChannel:
		badRequest(err, w)
		return
	case ids := <-proxyChannel:
		q.proxyAndSave(w, r, params, ids, engine)
		return
	case cache := <-cachedChannel:
		for _, i := range strings.Split(cache.headers, "\n") {
			h := strings.Split(i, "|")
			if len(h) >= 2 {
				w.Header().Set(h[0], h[1])
			}
		}

		w.Header().Set("x-Proxy", "served from cache")
		w.WriteHeader(cache.status)
		w.Write(cache.body)
		return
	}
}

func (q *QueryHandler) loadEngine(params *params.Params) (*engine.RuleEngine, error) {
	// With colleciton load rules for store
	rules, err := q.rules.Get(params.Collection)
	if err != nil {
		q.log.Warn("Failed to get rules")
		return nil, err
	}

	// Convert rules type
	engineRules := make([]engine.Rule, len(rules))
	for _, v := range rules {
		engineRules = append(engineRules, engine.Rule(v))
	}

	engine := &engine.RuleEngine{}

	// pass rules to engine
	engine.LoadRules(params.QueryURL, params.Collection, engineRules)

	return engine, nil
}

func (q *QueryHandler) checkResponses(
	params *params.Params, r *http.Request, engine *engine.RuleEngine,
	cachedChannel chan response, proxyChannel chan ids, errChannel chan error,
) {
	endpoint, err := q.endpoints.GetOrCreate(params.QueryURL.String(), r.Method, params.Collection)
	if err != nil {
		errChannel <- err
		return
	}

	res, err := q.responses.Get(
		params.QueryURL.String(),
		endpoint,
		r.Method,
	)
	if err != nil && err != fail.ErrNoData {
		errChannel <- err
		return
	}

	if err == fail.ErrNoData || res == nil {
		q.log.Info("No data found proxy request")
		proxyChannel <- ids{endpoint: endpoint, id: uuid.Nil}
		return

	}

	if !engine.HasExpired(res.DateTime) {
		q.log.Info("Returned saved response")
		cachedChannel <- response{
			headers: res.Headers,
			status:  res.Status,
			body:    res.Body,
		}
		return
	}

	q.log.Info("Response has expired - refresh data")
	proxyChannel <- ids{endpoint: endpoint, id: res.ID}
	return
}

func (q *QueryHandler) proxyAndSave(w http.ResponseWriter, r *http.Request, p *params.Params, ids ids, engine *engine.RuleEngine) {
	reverseProxy(w, r, p, func(re *http.Response) error {
		q.log.Info("Saving response for ", re.Request.URL)

		// Depulicate the body to reapply to response later
		buf, _ := ioutil.ReadAll(re.Body)
		body := new(bytes.Buffer)
		body.ReadFrom(ioutil.NopCloser(bytes.NewBuffer(buf)))

		headers := new(bytes.Buffer)
		for k, v := range re.Header {
			fmt.Fprintf(headers, "%s|%s\n", k, strings.Join(v, " "))
		}

		err := q.responses.Save(
			ids.id,
			re.Request.URL.String(),
			headers.String(),
			body.Bytes(),
			re.StatusCode,
			re.Request.Method,
			ids.endpoint,
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
		if engine.EnableCors() {
			corsHeaders(re.Header)
		}

		re.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
		return nil
	}, q.log)
}

func reverseProxy(
	w http.ResponseWriter, r *http.Request, p *params.Params, mr func(r *http.Response) error, logger *zap.SugaredLogger,
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
				logger.Warn("Internal Error on reverse Proxy - ", err)
				badRequest(fail.InternalError(err), w)
			}
		},
	}

	reverseProxy.ServeHTTP(w, r)
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
