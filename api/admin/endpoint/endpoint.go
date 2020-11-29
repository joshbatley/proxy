package endpoint

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/joshbatley/proxy/domain/endpoints"
	"github.com/joshbatley/proxy/internal/utils"
	"go.uber.org/zap"
)

type response struct {
	Count int        `json:"count"`
	Skip  int        `json:"skip"`
	Limit int        `json:"limit"`
	Data  []Endpoint `json:"data"`
}

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
	res := response{
		Count: len(es),
		Skip:  skip,
		Limit: limit,
		Data:  make([]Endpoint, 0),
	}

	for _, e := range es {
		res.Data = append(res.Data, Endpoint(e))
	}

	j, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json, text/plain, */*")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
