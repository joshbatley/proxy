package query

import (
	"github.com/google/uuid"
	"github.com/joshbatley/proxy/server/domain/collections"
	"github.com/joshbatley/proxy/server/domain/endpoints"
	"github.com/joshbatley/proxy/server/domain/responses"
	"github.com/joshbatley/proxy/server/domain/rules"
	"go.uber.org/zap"
)

type response struct {
	status  int
	body    []byte
	headers map[string]string
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
