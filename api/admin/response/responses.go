package response

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	domain "github.com/joshbatley/proxy/domain/responses"

	"github.com/joshbatley/proxy/internal/utils"
	"go.uber.org/zap"
)

type responses struct {
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
	responses *domain.Manager
	log       *zap.SugaredLogger
}

// NewHandler constructs a new QueryHandler
func NewHandler(
	responses *domain.Manager,
	log *zap.SugaredLogger,
) Handler {
	return Handler{
		responses,
		log,
	}
}

// Get -
func (h *Handler) Get(w http.ResponseWriter, re *http.Request) {
	p := re.URL.Query()
	skip, _ := strconv.Atoi(p.Get("skip"))
	limit, _ := strconv.Atoi(p.Get("limit"))
	endpoint := p.Get("endpoint")

	rs, err := h.responses.ListByEndpoint(endpoint, limit, skip)
	if err != nil {
		utils.BadRequest(err, w)
		return
	}

	data := []responses{}
	for _, r := range rs {
		data = append(data, responses{
			ID:       r.ID,
			Status:   r.Status,
			URL:      r.URL,
			Method:   r.Method,
			Headers:  r.Headers,
			Body:     string(r.Body),
			DateTime: r.DateTime,
		})
	}

	utils.Response(w, data, skip, limit)

}
