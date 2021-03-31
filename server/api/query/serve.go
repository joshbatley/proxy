package query

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
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
	res, err := q.QueryEngine(params, r.Method)
	if err != nil {
		utils.BadRequest(err, w)
		return
	}

	if res != nil && res.cache != nil {
		for k, v := range res.cache.headers {
			w.Header().Add(k, v)
		}
		w.WriteHeader(res.cache.status)
		w.Write(res.cache.body)
		return
	}

	if res != nil && res.proxyFunc != nil {
		reverseProxy(w, r, params, res.proxyFunc, q.log)
		return
	}

	// if sleepTime := engine.GetSleepTime(); sleepTime > 0 {
	//  q.log.Infof("Sleeping for %dms", sleepTime)
	// diff := time.Since(startTime).Milliseconds()
	// time.Sleep(time.Duration(sleepTime-diff) * time.Millisecond)
	// }

}

type QueryEngineResponse struct {
	cache     *response
	proxyFunc func(r *http.Response) error
}

func (q Handler) QueryEngine(p *params.Params, method string) (*QueryEngineResponse, error) {
	response := QueryEngineResponse{}

	// startTime := time.Now()
	engine, err := q.loadEngine(p)
	if err != nil {
		return nil, err
	}

	p.QueryURL = engine.Remapper()

	// Override cors request (Ignore if thirdparty API allow it or not)
	if method == http.MethodOptions && engine.EnableCors() {
		q.log.Info("Force CORS response")
		response.cache = cors()
		return &response, nil
	}

	// If the rule is to ignore
	if !engine.CheckStore() {
		response.proxyFunc = func(r *http.Response) error {
			return nil
		}
		return &response, nil
	}

	// Get Cache responses
	ids, cache, err := q.checkResponses(p, method, engine)

	// Ids exist so we can assume we need to reverseProxy
	// TODO clean up
	if ids.endpoint != uuid.Nil || ids.id != uuid.Nil {
		response.proxyFunc = q.proxyAndSave(ids, engine)
		return &response, nil
	}
	if cache != nil {
		q.log.Info("Returned saved response")
		body, err := encoder.Compress(cache.headers, cache.body)

		if err != nil {
			return nil, err
		}
		response.cache = cache
		response.cache.body = body
		return &response, nil
	}
	if err != nil {
		return nil, err
	}

	// Fallthrough
	return nil, nil
}

func cors() *response {
	return &response{
		body:   []byte{},
		status: 204,
		headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
			"Access-Control-Allow-Method": "*",
			"Access-Control-Allow-Header": "*",
		},
	}
}

func readHeaderString(hs string) map[string]string {
	n := map[string]string{}
	for _, i := range strings.Split(hs, "\n") {
		h := strings.Split(i, "|")
		if len(h) >= 2 {
			n[h[0]] = h[1]
		}
	}
	return n
}

func writeHeaderString(h http.Header) string {
	hs := new(bytes.Buffer)
	for k, v := range h {
		fmt.Fprintf(hs, "%s|%s\n", k, strings.Join(v, " "))
	}
	return hs.String()
}

func (q *Handler) proxyAndSave2(ids ids, engine *engine.RuleEngine) func(body []byte, headers http.Header, StatusCode int, method string, url string) error {
	return func(body []byte, headers http.Header, StatusCode int, method string, url string) error {

		return nil
	}
}

// body - bytes
// headers
// statusCode
// method
// url - string
//

func (q *Handler) proxyAndSave(ids ids, engine *engine.RuleEngine) func(r *http.Response) error {
	return func(r *http.Response) error {
		q.log.Info("Proxied and saving response for ", r.Request.URL)

		// Depulicate the body to reapply to response later
		buf, _ := ioutil.ReadAll(r.Body)
		body := new(bytes.Buffer)
		body.ReadFrom(ioutil.NopCloser(bytes.NewBuffer(buf)))
		r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

		content, err := encoder.Decompress(r.Header, body.Bytes())
		if err != nil {
			q.log.Info("Decoding body failed")
			content = body.Bytes()
		}

		h := new(bytes.Buffer)
		for k, v := range r.Header {
			fmt.Fprintf(h, "%s|%s\n", k, strings.Join(v, " "))
		}

		err = q.responses.Save(
			ids.id,
			r.Request.URL.String(),
			h.String(),
			content,
			r.StatusCode,
			r.Request.Method,
			ids.endpoint,
		)

		if err != nil {
			q.log.Error("Failed to save response")

			r.Header = http.Header{}
			r.Header.Set("Content-Type", "application/json, text/plain, */*")
			r.StatusCode = http.StatusBadRequest
			jsonString, _ := json.Marshal(fail.InternalError(err))
			r.Body = ioutil.NopCloser(bytes.NewBuffer(jsonString))
			return nil
		}

		// Apply headers to skip inbuild security
		if engine.EnableCors() {
			utils.Cors(r.Header)
		}

		return nil
	}
}

// func (q *Handler) proxyAndSave(w http.ResponseWriter, r *http.Request, p *params.Params, ids ids, engine *engine.RuleEngine) {
// reverseProxy(w, r, p, func(re *http.Response) error {
// q.log.Info("Saving response for ", re.Request.URL)

// // Depulicate the body to reapply to response later
// buf, _ := ioutil.ReadAll(re.Body)
// body := new(bytes.Buffer)
// body.ReadFrom(ioutil.NopCloser(bytes.NewBuffer(buf)))
// re.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

// content, err := encoder.Decompress(re.Header, body.Bytes())
// if err != nil {
// q.log.Info("Decoding body failed")
// content = body.Bytes()
// }

// h := new(bytes.Buffer)
// for k, v := range re.Header {
// fmt.Fprintf(h, "%s|%s\n", k, strings.Join(v, " "))
// }

// err = q.responses.Save(
// ids.id,
// re.Request.URL.String(),
// h.String(),
// content,
// re.StatusCode,
// re.Request.Method,
// ids.endpoint,
// )

// if err != nil {
// q.log.Error("Failed to save response")

// re.Header = http.Header{}
// re.Header.Set("Content-Type", "application/json, text/plain, */*")
// re.StatusCode = http.StatusBadRequest
// jsonString, _ := json.Marshal(fail.InternalError(err))
// re.Body = ioutil.NopCloser(bytes.NewBuffer(jsonString))
// return nil
// }

// // Apply headers to skip inbuild security
// if engine.EnableCors() {
// utils.Cors(re.Header)
// }

// return nil
// }, q.log)
// }

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

// return
// }
