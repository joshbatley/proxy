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
	params *params.Params, method string, engine *engine.RuleEngine,
) (ids, *response, error) {
	// init
	found := ids{endpoint: uuid.Nil, id: uuid.Nil}

	// Check for endpoint
	endpoint, err := q.endpoints.Get(params.QueryURL.String(), method, params.Collection)
	// Unexpected error
	if err != nil && err != fail.ErrNoData {
		return found, nil, err
	}
	// No endpoint
	// - Create new endpoints
	// - return new endpoint in ids or unexpected error
	if err == fail.ErrNoData {
		q.log.Info("New request, creating record and proxying", params.QueryURL)
		endpointID, err := q.endpoints.Save(params.QueryURL.String(), method, params.Collection)
		if err != nil {
			return found, nil, err
		}
		found.endpoint = endpointID
		return found, nil, nil
	}

	// Endpoint is found so get reponse
	res, err := q.responses.Get(
		params.QueryURL.String(),
		endpoint.ID,
		method,
		endpoint.Status,
	)

	// Unexpected Error
	if err != nil && err != fail.ErrNoData {
		return found, nil, err
	}

	// No Cache found this is probably a error
	if err == fail.ErrNoData || res == nil {
		q.log.Info("No response found, data is wrong")
		return found, nil, fail.ResponseMissing(err)
	}

	// Data has expired so return data
	if engine.HasExpired(res.DateTime) {
		q.log.Info("Response has expired - refresh data")
		return ids{
			endpoint: endpoint.ID, id: res.ID,
		}, nil, nil
	}

	// data is found so return cache
	return found, &response{
		headers: readHeaderString(res.Headers),
		status:  res.Status,
		body:    res.Body,
	}, nil
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
