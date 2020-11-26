package admin

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/joshbatley/proxy/internal/writers"
)

type responseResponse struct {
	Count int        `json:"count"`
	Skip  int        `json:"skip"`
	Limit int        `json:"limit"`
	Data  []response `json:"data"`
}

type response struct {
	ID       uuid.UUID `json:"id"`
	Status   int       `json:"status"`
	URL      string    `json:"url"`
	Method   string    `json:"method"`
	Headers  string    `json:"headers"`
	Body     string    `json:"body"`
	DateTime int64     `json:"datetime"`
}

func (h *Handler) response(w http.ResponseWriter, re *http.Request) {
	p := re.URL.Query()
	skip, _ := strconv.Atoi(p.Get("skip"))
	limit, _ := strconv.Atoi(p.Get("limit"))
	collectionID := p.Get("collection")

	rs, err := h.responses.ListByEndpoint(collectionID, limit, skip)
	if err != nil {
		log.Println(err)
		writers.BadRequest(err, w)
		return
	}
	var Responses []response

	for _, r := range rs {
		newHeader := make(http.Header, 0)
		for _, i := range strings.Split(r.Headers, "\n") {
			h := strings.Split(i, "|")
			if len(h) >= 2 {
				newHeader.Set(h[0], h[1])
			}
		}
		content, _ := writers.DecodeBody(newHeader, r.Body)
		Responses = append(Responses, response{
			ID:       r.ID,
			Status:   r.Status,
			URL:      r.URL,
			Method:   r.Method,
			Headers:  r.Headers,
			Body:     string(content),
			DateTime: r.DateTime,
		})
	}

	res := responseResponse{
		Count: len(rs),
		Skip:  skip,
		Limit: limit,
		Data:  Responses,
	}

	j, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json, text/plain, */*")

	// w.Header().Set("Content-Encoding", "br")
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}
