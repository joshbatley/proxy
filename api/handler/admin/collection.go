package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/joshbatley/proxy/domain/endpoints"
	"github.com/joshbatley/proxy/internal/writers"
)

type Collection struct {
	ID       string
	Name     string
	Endpoint *[]endpoints.Endpoint
}

func (h *Handler) collection(w http.ResponseWriter, r *http.Request) {
	cs, err := h.collections.List()
	if err != nil {
		writers.BadRequest(err, w)
		return
	}

	var res []Collection
	for _, v := range cs {
		d, err := h.endpoints.GetByID(v.ID)
		if err != nil {
			writers.BadRequest(err, w)
			return
		}
		res = append(res, Collection{
			ID:       strconv.FormatInt(v.ID, 10),
			Name:     v.Name,
			Endpoint: d,
		})
	}

	j, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json, text/plain, */*")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
