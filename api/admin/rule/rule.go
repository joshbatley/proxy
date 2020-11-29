package rule

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/joshbatley/proxy/domain/rules"
	"github.com/joshbatley/proxy/internal/utils"
	"go.uber.org/zap"
)

type response struct {
	Count int    `json:"count"`
	Skip  int    `json:"skip"`
	Limit int    `json:"limit"`
	Data  []rule `json:"data"`
}

type rule struct {
	Pattern      string `json:"pattern"`
	SaveResponse int    `json:"saveResponse"`
	ForceCors    int    `json:"sorceCors"`
	Expiry       int    `json:"expiry"`
	SkipOffline  int    `json:"skipOffline"`
	Delay        int    `json:"delay"`
	RemapRegex   string `json:"remapRegex"`
}

// Handler -
type Handler struct {
	rules *rules.Manager
	log   *zap.SugaredLogger
}

// NewHandler - Construct a new Handler
func NewHandler(rules *rules.Manager, log *zap.SugaredLogger,
) Handler {
	return Handler{
		rules,
		log,
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Query()
	skip, _ := strconv.Atoi(p.Get("skip"))
	limit, _ := strconv.Atoi(p.Get("limit"))
	collection := p.Get("collection")

	rs, err := h.rules.ListByCollectionID(collection, limit, skip)
	if err != nil {
		utils.BadRequest(err, w)
		return
	}

	res := response{
		Count: len(rs),
		Skip:  skip,
		Limit: limit,
		Data:  make([]rule, 0),
	}

	for _, r := range rs {
		res.Data = append(res.Data, rule{
			Pattern:      r.Pattern,
			SaveResponse: r.SaveResponse,
			ForceCors:    r.ForceCors,
			Expiry:       r.Expiry,
			SkipOffline:  r.SkipOffline,
			Delay:        r.Delay,
			RemapRegex:   r.RemapRegex,
		})
	}

	j, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json, text/plain, */*")
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}
