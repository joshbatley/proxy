package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joshbatley/proxy"
	"github.com/joshbatley/proxy/internal/engine"
	"github.com/joshbatley/proxy/internal/store"
	"github.com/joshbatley/proxy/internal/utils"
)

// QueryHandler Http handler for any query response
type QueryHandler struct {
	Store *store.Store
	Rules *engine.RuleEngine
}

// Serve Sets up all the logic for a reverse proxy and save and sends cached versions
func (q QueryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params, err := utils.ParseParams(mux.Vars(r), r.URL)
	if err != nil {
		badRequest(err, w)
		return
	}

	err = q.Rules.StartUp(params)
	if err != nil {
		badRequest(err, w)
		return
	}

	state := q.Rules.GetState()
	if err != nil {
		badRequest(err, w)
		return
	}

	if r.Method == http.MethodOptions && q.Rules.EnableCors() {
		corsHeaders(w.Header())
		w.WriteHeader(http.StatusNoContent)
		return
	}

	e, err := q.Store.GetOrAddEndpoint(params.QueryURL.String(), r.Method, params.Collection)
	if err != nil {
		badRequest(err, w)
		return
	}

	log.Println("We got here", e)

	switch state {
	case engine.StateSaving:
		d, err := q.Store.GetResponse(
			params.QueryURL.String(),
			e.ID,
			r.Method,
		)
		if err != nil {
			badRequest(err, w)
			return
		}
		if d != nil && !q.Rules.HasExpired(d.DateTime) {
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
	p *utils.Params,
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
			if state := q.Rules.GetState(); state == engine.StateSaving {
				// Apply headers to skip inbuild security
				if q.Rules.EnableCors() {
					corsHeaders(re.Header)
				}

				// Depulicate the body to reapply to response later
				buf, _ := ioutil.ReadAll(re.Body)
				err := q.Store.SaveResponse(
					proxy.NewResponse(
						re.Request.URL,
						ioutil.NopCloser(bytes.NewBuffer(buf)),
						re.Header,
						re.StatusCode,
						re.Request.Method,
						p.Collection,
					),
				)
				if err != nil {
					badResponse(proxy.InternalError(err), re)
					return nil
				}
				re.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
				return nil
			}
			return nil
		},
	}

	reverseProxy.ServeHTTP(w, r)
}

func sendResponse(d *proxy.ResponseRow, w http.ResponseWriter) {
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
		jsonString, _ = json.Marshal(proxy.InternalError(err))
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(jsonString))
}

func badRequest(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json, text/plain, */*")
	w.WriteHeader(http.StatusBadRequest)
	jsonString, _ := json.Marshal(err)
	if len(jsonString) == 2 {
		jsonString, _ = json.Marshal(proxy.InternalError(err))
	}
	w.Write(jsonString)
}

func corsHeaders(h http.Header) {
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Access-Control-Allow-Methods", "*")
	h.Set("Access-Control-Allow-Headers", "*")
}
