package query

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joshbatley/proxy/internal/engine"
	"github.com/joshbatley/proxy/internal/fail"
	"github.com/joshbatley/proxy/internal/params"
	"github.com/joshbatley/proxy/internal/writers"
)

type ModifyRsponse func(re *http.Response) error

// Serve Sets up all the logic for a reverse proxy and save and sends cached versions
func (q Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	params, err := params.Parse(mux.Vars(r), r.URL)
	if err != nil {
		q.log.Error("Param parse failed")
		writers.BadRequest(err, w)
		return
	}

	engine, err := q.loadEngine(params)
	if err != nil {
		q.log.Warn("Failed to load rules")
		writers.BadRequest(err, w)
		return
	}
	params.QueryURL = engine.Remapper()

	// Check if method is OPTIONS and if Engine need to override
	if r.Method == http.MethodOptions && engine.EnableCors() {
		q.log.Info("Cors request")
		writers.CorsHeaders(w.Header())
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Skip Store check, and straight proxy
	if !engine.CheckStore() {
		q.log.Info("Skipping cache check and proxing: ", params.QueryURL)
		reverseProxy(w, r, params, func(re *http.Response) error {
			if engine.EnableCors() {
				writers.CorsHeaders(re.Header)
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
		writers.BadRequest(err, w)
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
		if sleepTime := engine.GetSleepTime(); sleepTime > 0 {
			q.log.Infof("Sleeping for %dms", sleepTime)
			diff := time.Since(startTime).Milliseconds()
			time.Sleep(time.Duration(sleepTime-diff) * time.Millisecond)
		}
		q.log.Info("Returned saved response")

		w.Header().Set("x-Proxy", "served from cache")
		w.WriteHeader(cache.status)
		w.Write(cache.body)
		return
	}
}

func (q *Handler) loadEngine(params *params.Params) (*engine.RuleEngine, error) {
	// Check the collection exist (if not default)
	col, err := q.collections.Get(params.Collection)
	if err == fail.ErrNoData {
		q.log.Warn("No collection found")
		return nil, fail.MissingColErr(err)
	}
	var urls []string
	if col != nil {
		urls = strings.Split(col.HealthCheckURLs.String, ",")
	}

	// With colleciton load rules for store
	rules, err := q.rules.Get(params.Collection)
	if err != nil {
		q.log.Warn("Failed to get rules", err)
		return nil, err
	}

	// Convert rules type
	engineRules := make([]engine.Rule, len(rules))
	for _, v := range rules {
		engineRules = append(engineRules, engine.Rule(v))
	}

	engine := &engine.RuleEngine{}

	// pass rules to engine
	engine.LoadRules(params.QueryURL, params.Collection, engineRules, urls)

	return engine, nil
}

func (q *Handler) checkResponses(
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

func (q *Handler) proxyAndSave(w http.ResponseWriter, r *http.Request, p *params.Params, ids ids, engine *engine.RuleEngine) {
	reverseProxy(w, r, p, func(re *http.Response) error {
		q.log.Info("Saving response for ", re.Request.URL)

		// Depulicate the body to reapply to response later
		buf, _ := ioutil.ReadAll(re.Body)
		body := new(bytes.Buffer)
		body.ReadFrom(ioutil.NopCloser(bytes.NewBuffer(buf)))
		re.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

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
			writers.CorsHeaders(re.Header)
		}

		return nil
	}, q.log)
}
