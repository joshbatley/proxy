package collection

import (
	"net/http"
	"strconv"

	"github.com/joshbatley/proxy/server/api/admin/endpoint"
	"github.com/joshbatley/proxy/server/domain/collections"
	"github.com/joshbatley/proxy/server/domain/endpoints"
	"github.com/joshbatley/proxy/server/internal/utils"
	"go.uber.org/zap"
)

type collection struct {
	ID        string              `json:"id"`
	Name      string              `json:"name"`
	Endpoints []endpoint.Endpoint `json:"endpoints,omitempty"`
}

// Handler requires collections, endpoints and logger
type Handler struct {
	collections *collections.Manager
	endpoints   *endpoints.Manager
	log         *zap.SugaredLogger
}

// NewHandler construct a new Handler
func NewHandler(collections *collections.Manager, endpoints *endpoints.Manager, log *zap.SugaredLogger,
) Handler {
	return Handler{
		collections,
		endpoints,
		log,
	}
}

// GetCollections returns all collection with pagination
func (h *Handler) GetCollections(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Query()
	skip, _ := strconv.Atoi(p.Get("skip"))
	limit, _ := strconv.Atoi(p.Get("limit"))

	cs, err := h.collections.List(limit, skip)
	if err != nil {
		utils.BadRequest(err, w)
		return
	}

	data := []collection{}
	for _, v := range cs {
		d, err := h.endpoints.GetByCollectionID(v.ID)
		if err != nil {
			utils.BadRequest(err, w)
			return
		}

		var endpoints []endpoint.Endpoint
		for _, e := range *d {
			endpoints = append(endpoints, endpoint.Endpoint(e))
		}

		data = append(data, collection{
			ID:        strconv.FormatInt(v.ID, 10),
			Name:      v.Name,
			Endpoints: endpoints,
		})
	}

	utils.PaginatedWrap(w, data, limit, skip)
}
