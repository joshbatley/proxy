package query

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"

	"github.com/joshbatley/proxy/internal/connection"
	"github.com/joshbatley/proxy/internal/fail"
	"github.com/joshbatley/proxy/internal/params"
	"go.uber.org/zap"
)

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
			if connection.IsOffline(nil) {
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
