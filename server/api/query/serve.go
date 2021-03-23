package query

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joshbatley/proxy/server/internal/encoder"
	"github.com/joshbatley/proxy/server/internal/engine"
	"github.com/joshbatley/proxy/server/internal/fail"
	"github.com/joshbatley/proxy/server/internal/params"
	"github.com/joshbatley/proxy/server/internal/utils"
)

// ModifyResponse required stuct
type ModifyResponse func(re *http.Response) error

func (q Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	params, err := params.Parse(mux.Vars(r), r.URL)
	// startTime := time.Now()
	// params, err := params.Parse(mux.Vars(r), r.URL)
	// Todo seperate the query logic and the http handler
	// ok, cache, err := QyeryEngine()
	// ok = just proxy
	// cache = datafound
	// err = badrequest

	// To think about
	// CORS request
	// ignore saving
	// proxy and save
	res, err := q.QueryEngine(params, r)
	if err != nil {
		utils.BadRequest(err, w)
		return
	}
	if res != nil && res.cache != nil {
		for _, i := range strings.Split(res.cache.headers, "\n") {
			h := strings.Split(i, "|")
			if len(h) >= 2 {
				w.Header().Set(h[0], h[1])
			}
		}
		w.WriteHeader(res.cache.status)
		w.Write(res.cache.body)
		return
	}

	// if sleepTime := engine.GetSleepTime(); sleepTime > 0 {
	//  q.log.Infof("Sleeping for %dms", sleepTime)
	// diff := time.Since(startTime).Milliseconds()
	// time.Sleep(time.Duration(sleepTime-diff) * time.Millisecond)
	// }

	// take respone
	// desctruct header
	// apply to writer
}

type QueryEngineResponse struct {
	cache     *response
	proxyFunc func(r *http.Response) error
}

func (q Handler) QueryEngine(p *params.Params, r *http.Request) (*QueryEngineResponse, error) {
	response := QueryEngineResponse{}

	// startTime := time.Now()
	engine, err := q.loadEngine(p)
	if err != nil {
		return nil, err
	}

	p.QueryURL = engine.Remapper()

	if r.Method == http.MethodOptions && engine.EnableCors() {
		q.log.Info("Force CORS response")
		response.cache = cors()
		return &response, nil
	}

	ids, cache, err := q.checkResponses(p, r, engine)
	log.Println(ids)
	if cache != nil {

		q.log.Info("Returned saved response")
		body, err := encoder.Compress(cache.headers, cache.body)

		if err != nil {
			return nil, err
		}
		response.cache = cache
		response.cache.headers = cache.headers + "x-proxy|served from cache"
		response.cache.body = body
		return &response, nil
	}
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func cors() *response {
	return &response{
		headers: "Access-Control-Allow-Origin|*\nAccess-Control-Allow-Methods|*\nAccess-Control-Allow-Headers|*",
		body:    []byte{},
		status:  204,
	}
}

// func Cors(h http.Header) {
// h.Set("Access-Control-Allow-Origin", "*")
// h.Set("Access-Control-Allow-Methods", "*")
// h.Set("Access-Control-Allow-Headers", "*")

func (q *Handler) proxyAndSave(w http.ResponseWriter, r *http.Request, p *params.Params, ids ids, engine *engine.RuleEngine) {
	reverseProxy(w, r, p, func(re *http.Response) error {
		q.log.Info("Saving response for ", re.Request.URL)

		// Depulicate the body to reapply to response later
		buf, _ := ioutil.ReadAll(re.Body)
		body := new(bytes.Buffer)
		body.ReadFrom(ioutil.NopCloser(bytes.NewBuffer(buf)))
		re.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

		content, err := encoder.Decompress(re.Header, body.Bytes())
		if err != nil {
			q.log.Info("Decoding body failed")
			content = body.Bytes()
		}

		h := new(bytes.Buffer)
		for k, v := range re.Header {
			fmt.Fprintf(h, "%s|%s\n", k, strings.Join(v, " "))
		}

		err = q.responses.Save(
			ids.id,
			re.Request.URL.String(),
			h.String(),
			content,
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
			utils.Cors(re.Header)
		}

		return nil
	}, q.log)
}

// ServeHTTP sets up all the logic for a reverse proxy and save and sends cached versions
// func (q Handler) QueryEngine2(w http.ResponseWriter, r *http.Request) {
// startTime := time.Now()
// params, err := params.Parse(mux.Vars(r), r.URL)
// if err != nil {
// q.log.Error("Param parse failed")
// utils.BadRequest(err, w)
// return
// }
// engine, err := q.loadEngine(params)
// if err != nil {
// q.log.Warn("Failed to load rules")
// utils.BadRequest(err, w)
// return
// }
// params.QueryURL = engine.Remapper()

// // Check if method is OPTIONS and if Engine need to override
// if r.Method == http.MethodOptions && engine.EnableCors() {
// q.log.Info("Cors request")
// utils.Cors(w.Header())
// w.WriteHeader(http.StatusNoContent)
// return
// }

// // Skip Store check, and straight proxy
// if !engine.CheckStore() {
// q.log.Info("Skipping cache check and proxing: ", params.QueryURL)
// reverseProxy(w, r, params, func(re *http.Response) error {
// if engine.EnableCors() {
// utils.Cors(re.Header)
// }
// return nil
// }, q.log)
// return
// }

// ids, cache, err := q.checkResponses(params, r, engine)

// if err != nil {
// utils.BadRequest(err, w)
// return
// }

// if ids.endpoint != uuid.Nil || ids.id != uuid.Nil {
// q.proxyAndSave(w, r, params, ids, engine)
// return
// }

// //	q.sendCachedResponse(w, cache, engine, startTime)
// return
// }
