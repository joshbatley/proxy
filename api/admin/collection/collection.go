package collection

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/joshbatley/proxy/domain/collections"
	"github.com/joshbatley/proxy/domain/endpoints"
	"github.com/joshbatley/proxy/internal/utils"
	"go.uber.org/zap"
)

type response struct {
	Count int          `json:"count"`
	Skip  int          `json:"skip"`
	Limit int          `json:"limit"`
	Data  []collection `json:"data"`
}

type collection struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Endpoints []endpoint `json:"endpoints,omitempty"`
}

type endpoint struct {
	ID     uuid.UUID `json:"id"`
	Status int       `json:"status"`
	Method string    `json:"method"`
	URL    string    `json:"url"`
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

	res := response{
		Count: len(cs),
		Skip:  skip,
		Limit: limit,
		Data:  make([]collection, 0),
	}

	for _, v := range cs {
		d, err := h.endpoints.GetByID(v.ID)
		if err != nil {
			utils.BadRequest(err, w)
			return
		}

		var Endpoints []endpoint
		for _, e := range *d {
			Endpoints = append(Endpoints, endpoint(e))
		}

		res.Data = append(res.Data, collection{
			ID:        strconv.FormatInt(v.ID, 10),
			Name:      v.Name,
			Endpoints: Endpoints,
		})
	}

	j, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json, text/plain, */*")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
