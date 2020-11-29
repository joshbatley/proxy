package admin

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
	r.PathPrefix("/collections").Methods("GET").Handler(addGzip(h.collection))
	r.PathPrefix("/endpoints").Methods("GET").Handler(addGzip(h.endpoint))
	r.PathPrefix("/responses").Methods("GET").Handler(addGzip(h.response))
	r.NotFoundHandler = http.HandlerFunc(utils.NotFound)
}

func addGzip(h http.HandlerFunc) http.Handler {
	return handlers.CompressHandler(http.HandlerFunc(h))
}
