package response

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/joshbatley/proxy/domain/responses"
	"github.com/joshbatley/proxy/internal/utils"
	"go.uber.org/zap"
)

type response struct {
	Count int    `json:"count"`
	Skip  int    `json:"skip"`
	Limit int    `json:"limit"`
	Data  []data `json:"data"`
}

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
	responses *responses.Manager
	log       *zap.SugaredLogger
}

// NewHandler constructs a new QueryHandler
func NewHandler(
	responses *responses.Manager,
	log *zap.SugaredLogger,
) Handler {
	return Handler{
		responses,
		log,
	}
}

// Get
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

	res := response{
		Count: len(rs),
		Skip:  skip,
		Limit: limit,
		Data:  make([]data, 0),
	}

	for _, r := range rs {
		res.Data = append(res.Data, data{
			ID:       r.ID,
			Status:   r.Status,
			URL:      r.URL,
			Method:   r.Method,
			Headers:  r.Headers,
			Body:     string(r.Body),
			DateTime: r.DateTime,
		})
	}

	j, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json, text/plain, */*")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
