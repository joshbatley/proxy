package endpoint

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joshbatley/proxy/server/domain/endpoints"
	"github.com/joshbatley/proxy/server/internal/utils"
	"go.uber.org/zap"
)

// Endpoint json response type
type Endpoint struct {
	ID           uuid.UUID `json:"id"`
	Status       int       `json:"status"`
	Method       string    `json:"method"`
	URL          string    `json:"url"`
	CollectionID string    `json:"collectionId"`
}

// Handler requires endpoint and logger
type Handler struct {
	endpoints *endpoints.Manager
	log       *zap.SugaredLogger
}

// NewHandler construct a new Handler
func NewHandler(endpoints *endpoints.Manager, log *zap.SugaredLogger,
) Handler {
	return Handler{
		endpoints,
		log,
	}
}

// Get all endpoints
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Query()
	skip, _ := strconv.Atoi(p.Get("skip"))
	limit, _ := strconv.Atoi(p.Get("limit"))

	es, err := h.endpoints.List(limit, skip)
	if err != nil {
		utils.BadRequest(err, w)
		return
	}

	data := []Endpoint{}
	for _, e := range es {
		data = append(data, Endpoint(e))
	}

	utils.PaginatedWrap(w, data, limit, skip)
}

// GetByID returns all endpoints by ID
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	es, err := h.endpoints.GetByID(id)
	if err != nil {
		utils.BadRequest(err, w)
		return
	}

	data := Endpoint(*es)
	j, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json, text/plain, */*")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
