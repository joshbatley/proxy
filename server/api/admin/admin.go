package admin

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joshbatley/proxy/server/api/admin/collection"
	"github.com/joshbatley/proxy/server/api/admin/endpoint"
	"github.com/joshbatley/proxy/server/api/admin/response"
	"github.com/joshbatley/proxy/server/api/admin/rule"
	"github.com/joshbatley/proxy/server/domain/collections"
	"github.com/joshbatley/proxy/server/domain/endpoints"
	"github.com/joshbatley/proxy/server/domain/responses"
	"github.com/joshbatley/proxy/server/domain/rules"
	"github.com/joshbatley/proxy/server/internal/utils"
	"go.uber.org/zap"
)

// Handler Http handler for any query response
type Handler struct {
	collections *collections.Manager
	endpoints   *endpoints.Manager
	responses   *responses.Manager
	rules       *rules.Manager
	log         *zap.SugaredLogger
}

// NewHandler constructs a new QueryHandler
func NewHandler(
	collections *collections.Manager,
	endpoints *endpoints.Manager,
	responses *responses.Manager,
	rules *rules.Manager,
	log *zap.SugaredLogger,
) Handler {
	return Handler{
		collections: collections,
		endpoints:   endpoints,
		responses:   responses,
		rules:       rules,
		log:         log,
	}
}

// Router returns a new admin router
func (h Handler) Router(r *mux.Router) {
	cols := collection.NewHandler(h.collections, h.endpoints, h.log)
	r.PathPrefix("/collections").Methods("GET").Handler(addGzip(cols.GetCollections))

	ends := endpoint.NewHandler(h.endpoints, h.log)
	r.PathPrefix("/endpoints/{id}").Methods("GET").Handler(addGzip(ends.GetByID))
	r.PathPrefix("/endpoints").Methods("GET").Handler(addGzip(ends.Get))

	res := response.NewHandler(h.responses, h.log)
	r.PathPrefix("/responses/{id}").Methods("GET").Handler(addGzip(res.Get))

	rule := rule.NewHandler(h.rules, h.log)
	r.PathPrefix("/rules").Methods("GET").Handler(addGzip(rule.Get))

	r.NotFoundHandler = http.HandlerFunc(utils.NotFound)
	r.Use(handlers.CORS())
}

func addGzip(h http.HandlerFunc) http.Handler {
	return handlers.CompressHandler(http.HandlerFunc(h))
}
