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

func NewQueryHandler(
	col *collections.Manager,
	end *endpoints.Manager,
	res *responses.Manager,
	rul *rules.Manager,
	eng *engine.RuleEngine) QueryHandler {
	return QueryHandler{
		collections: col,
		endpoints:   end,
		responses:   res,
		rules:       rul,
		engine:      eng,
	}
}

// Serve Sets up all the logic for a reverse proxy and save and sends cached versions
func (q QueryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params, err := params.Parse(mux.Vars(r), r.URL)
	if err != nil {
		badRequest(err, w)
		return
	}

	err = q.engine.LoadRules(params)
	if err != nil {
		badRequest(err, w)
		return
	}

	if r.Method == http.MethodOptions && q.engine.EnableCors() {
		corsHeaders(w.Header())
		w.WriteHeader(http.StatusNoContent)
		return
	}

	state := q.engine.GetState()

	e, err := q.endpoints.GetOrSave(params.QueryURL.String(), r.Method, params.Collection)
	if err != nil {
		badRequest(err, w)
		return
	}

	switch state {
	case engine.StateSaving:
		d, err := q.responses.Get(
			params.QueryURL.String(),
			e.ID,
			r.Method,
		)
		if err != nil {
			badRequest(err, w)
			return
		}
		if d != nil && !q.engine.HasExpired(d.DateTime) {
			sendResponse(d, w)
			return
		}
		fallthrough
	default:
		q.reverseProxy(w, r, params)
	}
}

func (q QueryHandler) reverseProxy(
	w http.ResponseWriter,
	r *http.Request,
	p *params.Params,
) {
	reverseProxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Header.Del("Origin")
			req.Header.Del("Referer")
			req.URL.Scheme = p.QueryURL.Scheme
			req.URL.Host = p.QueryURL.Host
			req.URL.Path = p.QueryURL.Path
			req.Host = p.QueryURL.Host
			req.URL.RawQuery = p.QueryURL.RawQuery
		},
		ModifyResponse: func(re *http.Response) error {
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
		},
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
