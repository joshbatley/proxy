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
	"github.com/joshbatley/proxy/server/internal/encoder"
	"github.com/joshbatley/proxy/server/internal/engine"
	"github.com/joshbatley/proxy/server/internal/fail"
	"github.com/joshbatley/proxy/server/internal/params"
	"github.com/joshbatley/proxy/server/internal/utils"
)

// ModifyRsponse -
type ModifyRsponse func(re *http.Response) error

// Serve Sets up all the logic for a reverse proxy and save and sends cached versions
func (q Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	params, err := params.Parse(mux.Vars(r), r.URL)
	if err != nil {
		q.log.Error("Param parse failed")
		utils.BadRequest(err, w)
		return
	}

	engine, err := q.loadEngine(params)
	if err != nil {
		q.log.Warn("Failed to load rules")
		utils.BadRequest(err, w)
		return
	}
	params.QueryURL = engine.Remapper()

	// Check if method is OPTIONS and if Engine need to override
	if r.Method == http.MethodOptions && engine.EnableCors() {
		q.log.Info("Cors request")
		utils.Cors(w.Header())
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Skip Store check, and straight proxy
	if !engine.CheckStore() {
		q.log.Info("Skipping cache check and proxing: ", params.QueryURL)
		reverseProxy(w, r, params, func(re *http.Response) error {
			if engine.EnableCors() {
				utils.Cors(re.Header)
			}
			return nil
		}, q.log)
		return
	}

	ids, cache, err := q.checkResponses(params, r, engine)

	if err != nil {
		utils.BadRequest(err, w)
		return
	}

	if ids.endpoint != uuid.Nil || ids.id != uuid.Nil {
		q.proxyAndSave(w, r, params, ids, engine)
		return
	}

	q.sendCachedResponse(w, cache, engine, startTime)
	return
}

func (q *Handler) sendCachedResponse(w http.ResponseWriter, cache *response, engine *engine.RuleEngine, startTime time.Time) {
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
	body, err := encoder.Compress(w.Header(), cache.body)
	if err != nil {
		utils.BadRequest(err, w)
		return
	}
	w.Header().Set("x-Proxy", "served from cache")
	w.WriteHeader(cache.status)
	w.Write(body)
}

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
