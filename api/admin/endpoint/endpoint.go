package endpoint

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joshbatley/proxy/domain/endpoints"
	"github.com/joshbatley/proxy/internal/utils"
	"go.uber.org/zap"
)

type Endpoint struct {
	ID     uuid.UUID `json:"id"`
	Status int       `json:"status"`
	Method string    `json:"method"`
	URL    string    `json:"url"`
}

// Handler -
type Handler struct {
	endpoints *endpoints.Manager
	log       *zap.SugaredLogger
}

// NewHandler - Construct a new Handler
func NewHandler(endpoints *endpoints.Manager, log *zap.SugaredLogger,
) Handler {
	return Handler{
		endpoints,
		log,
	}
}

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

	utils.Response(w, data, skip, limit)
}

func (h *Handler) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	es, err := h.endpoints.GetByID(id)
	if err != nil {
		utils.BadRequest(err, w)
		return
	}

	data := []Endpoint{}
	for _, e := range es {
		data = append(data, Endpoint(e))
	}

	j, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json, text/plain, */*")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
