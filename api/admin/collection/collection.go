package collection

import (
	"net/http"
	"strconv"

	"github.com/joshbatley/proxy/api/admin/endpoint"
	"github.com/joshbatley/proxy/domain/collections"
	"github.com/joshbatley/proxy/domain/endpoints"
	"github.com/joshbatley/proxy/internal/utils"
	"go.uber.org/zap"
)

type collection struct {
	ID        string              `json:"id"`
	Name      string              `json:"name"`
	Endpoints []endpoint.Endpoint `json:"endpoints,omitempty"`
}

// Handler -
type Handler struct {
	collections *collections.Manager
	endpoints   *endpoints.Manager
	log         *zap.SugaredLogger
}

// NewHandler - Construct a new Handler
func NewHandler(collections *collections.Manager, endpoints *endpoints.Manager, log *zap.SugaredLogger,
) Handler {
	return Handler{
		collections,
		endpoints,
		log,
	}
}

func (h *Handler) Selector(w http.ResponseWriter, r *http.Request) {
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
		d, err := h.endpoints.GetByColID(v.ID)
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

	utils.Response(w, data, limit, skip)
}
