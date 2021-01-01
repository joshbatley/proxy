package response

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joshbatley/proxy/server/domain/responses"
	"github.com/joshbatley/proxy/server/internal/utils"
	"go.uber.org/zap"
)

type data struct {
	ID       uuid.UUID `json:"id"`
	Status   int       `json:"status"`
	URL      string    `json:"url"`
	Method   string    `json:"method,omitempty"`
	Headers  string    `json:"headers,omitempty"`
	Body     string    `json:"body,omitempty"`
	DateTime int64     `json:"datetime,omitempty"`
}

// Handler Http handler for any query response
type Handler struct {
	repo *responses.Manager
	log  *zap.SugaredLogger
}

// NewHandler constructs a new QueryHandler
func NewHandler(
	repo *responses.Manager,
	log *zap.SugaredLogger,
) Handler {
	return Handler{
		repo,
		log,
	}
}

// Get all by response by endpoint ID
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	endpoint, err := uuid.Parse(vars["id"])

	p := r.URL.Query()
	skip, _ := strconv.Atoi(p.Get("skip"))
	limit, _ := strconv.Atoi(p.Get("limit"))

	rs, err := h.repo.ListByEndpoint(endpoint, limit, skip)
	if err != nil {
		utils.BadRequest(err, w)
		return
	}

	d := []data{}
	for _, res := range rs {
		d = append(d, data{
			ID:     res.ID,
			Status: res.Status,
			URL:    res.URL,
		})
	}

	utils.PaginatedWrap(w, d, limit, skip)

}
