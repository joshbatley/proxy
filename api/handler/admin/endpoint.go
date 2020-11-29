package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/joshbatley/proxy/internal/utils"
)

type endpointResponse struct {
	Count int        `json:"count"`
	Skip  int        `json:"skip"`
	Limit int        `json:"limit"`
	Data  []endpoint `json:"data"`
}

func (h *Handler) endpoint(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Query()
	skip, _ := strconv.Atoi(p.Get("skip"))
	limit, _ := strconv.Atoi(p.Get("limit"))

	es, err := h.endpoints.List(limit, skip)
	if err != nil {
		utils.BadRequest(err, w)
		return
	}

	var Endpoints []endpoint
	for _, e := range es {
		Endpoints = append(Endpoints, endpoint(e))
	}

	res := endpointResponse{
		Count: len(es),
		Skip:  skip,
		Limit: limit,
		Data:  Endpoints,
	}

	j, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json, text/plain, */*")
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}
