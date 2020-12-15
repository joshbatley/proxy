package response

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/joshbatley/proxy/domain/responses"
	"github.com/joshbatley/proxy/internal/utils"
	"go.uber.org/zap"
)

type data struct {
	ID       uuid.UUID `json:"id"`
	Status   int       `json:"status"`
	URL      string    `json:"url"`
	Method   string    `json:"method"`
	Headers  string    `json:"headers"`
	Body     string    `json:"body"`
	DateTime int64     `json:"datetime"`
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

// Get -
func (h *Handler) Get(w http.ResponseWriter, re *http.Request) {
	p := re.URL.Query()
	skip, _ := strconv.Atoi(p.Get("skip"))
	limit, _ := strconv.Atoi(p.Get("limit"))
	endpoint := p.Get("endpoint")

	rs, err := h.repo.ListByEndpoint(endpoint, limit, skip)
	if err != nil {
		utils.BadRequest(err, w)
		return
	}

	d := []data{}
	for _, r := range rs {
		d = append(d, data{
			ID:       r.ID,
			Status:   r.Status,
			URL:      r.URL,
			Method:   r.Method,
			Headers:  r.Headers,
			Body:     string(r.Body),
			DateTime: r.DateTime,
		})
	}

	utils.Response(w, d, skip, limit)

}
