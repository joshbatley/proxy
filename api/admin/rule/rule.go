package rule

import (
	"net/http"
	"strconv"

	"github.com/joshbatley/proxy/domain/rules"
	"github.com/joshbatley/proxy/internal/utils"
	"go.uber.org/zap"
)

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

	data := []rule{}
	for _, r := range rs {
		data = append(data, rule{
			Pattern:      r.Pattern,
			SaveResponse: r.SaveResponse,
			ForceCors:    r.ForceCors,
			Expiry:       r.Expiry,
			SkipOffline:  r.SkipOffline,
			Delay:        r.Delay,
			RemapRegex:   r.RemapRegex,
		})
	}
	utils.Response(w, data, limit, skip)
}
