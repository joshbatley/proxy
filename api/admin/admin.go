package admin

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joshbatley/proxy/api/admin/collection"
	"github.com/joshbatley/proxy/api/admin/endpoint"
	"github.com/joshbatley/proxy/api/admin/response"
	"github.com/joshbatley/proxy/api/admin/rule"
	"github.com/joshbatley/proxy/domain/collections"
	"github.com/joshbatley/proxy/domain/endpoints"
	"github.com/joshbatley/proxy/domain/responses"
	"github.com/joshbatley/proxy/domain/rules"
	"github.com/joshbatley/proxy/internal/utils"
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

// Router -
func (h Handler) Router(r *mux.Router) {
	cols := collection.NewHandler(h.collections, h.endpoints, h.log)
	ends := endpoint.NewHandler(h.endpoints, h.log)
	res := response.NewHandler(h.responses, h.log)
	rule := rule.NewHandler(h.rules, h.log)

	r.PathPrefix("/collections/selector").Methods("GET").Handler(addGzip(cols.Selector))
	r.PathPrefix("/endpoints").Methods("GET").Handler(addGzip(ends.Get))
	r.PathPrefix("/responses").Methods("GET").Handler(addGzip(res.Get))
	r.PathPrefix("/rules").Methods("GET").Handler(addGzip(rule.Get))

	r.NotFoundHandler = http.HandlerFunc(utils.NotFound)
	r.Use(handlers.CORS())
}

func addGzip(h http.HandlerFunc) http.Handler {
	return handlers.CompressHandler(http.HandlerFunc(h))
}
