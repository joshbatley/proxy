package admin

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joshbatley/proxy/domain/collections"
	"github.com/joshbatley/proxy/domain/endpoints"
	"github.com/joshbatley/proxy/domain/responses"
	"github.com/joshbatley/proxy/domain/rules"
	"github.com/joshbatley/proxy/internal/writers"
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
	r.PathPrefix("/collections").Methods("GET").HandlerFunc(h.collection)
	r.NotFoundHandler = http.HandlerFunc(writers.ReturnNotFound)

}
