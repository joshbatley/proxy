package query

import (
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/google/uuid"
	"github.com/joshbatley/proxy/server/internal/connection"
	"github.com/joshbatley/proxy/server/internal/engine"
	"github.com/joshbatley/proxy/server/internal/fail"
	"github.com/joshbatley/proxy/server/internal/params"
	"github.com/joshbatley/proxy/server/internal/utils"
	"go.uber.org/zap"
)

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
	rules, err := q.rules.GetByCollectionID(params.Collection)
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
) (ids, *response, error) {
	found := ids{endpoint: uuid.Nil, id: uuid.Nil}

	endpoint, err := q.endpoints.Get(params.QueryURL.String(), r.Method, params.Collection)
	if err != nil && err != fail.ErrNoData {
		return found, nil, err
	}
	if err == fail.ErrNoData {
		q.log.Info("New request, creating record and proxying", params.QueryURL)
		endpointID, err := q.endpoints.Save(params.QueryURL.String(), r.Method, params.Collection)
		if err != nil {
			return found, nil, err
		}
		found.endpoint = endpointID
		return found, nil, nil
	}

	res, err := q.responses.Get(
		params.QueryURL.String(),
		endpoint.ID,
		r.Method,
		endpoint.Status,
	)

	if err != nil && err != fail.ErrNoData {
		return found, nil, err
	}

	if err == fail.ErrNoData || res == nil {
		q.log.Info("No response found, data is wrong")
		return found, nil, fail.ResponseMissing(err)
	}

	if !engine.HasExpired(res.DateTime) {
		return found, &response{
			headers: res.Headers,
			status:  res.Status,
			body:    res.Body,
		}, nil

	}

	q.log.Info("Response has expired - refresh data")
	return ids{
		endpoint: endpoint.ID, id: res.ID,
	}, nil, nil
}

func reverseProxy(
	w http.ResponseWriter, r *http.Request, p *params.Params, mr ModifyResponse, logger *zap.SugaredLogger,
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
				utils.BadRequest(fail.OfflineError(err), w)
			} else {
				logger.Warn("Internal Error on reverse Proxy - ", err)
				utils.BadRequest(fail.InternalError(err), w)
			}
		},
	}

	reverseProxy.ServeHTTP(w, r)
}
