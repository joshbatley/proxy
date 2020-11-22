package query

import (
	"github.com/google/uuid"
	"github.com/joshbatley/proxy/domain/collections"
	"github.com/joshbatley/proxy/domain/endpoints"
	"github.com/joshbatley/proxy/domain/responses"
	"github.com/joshbatley/proxy/domain/rules"
	"go.uber.org/zap"
)

type response struct {
	headers string
	status  int
	body    []byte
}

type ids struct {
	endpoint uuid.UUID
	id       uuid.UUID
}

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
