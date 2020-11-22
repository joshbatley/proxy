package writers

import (
	"encoding/json"
	"net/http"

	"github.com/joshbatley/proxy/internal/fail"
)

// ReturnNotFound -
func ReturnNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json, text/plain, */*")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{}`))
}

// BadRequest -
func BadRequest(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json, text/plain, */*")
	w.WriteHeader(http.StatusBadRequest)
	jsonString, _ := json.Marshal(err)

	if len(jsonString) == 2 {
		jsonString, _ = json.Marshal(fail.InternalError(err))
	}

	w.Write(jsonString)
}

// CorsHeaders -
func CorsHeaders(h http.Header) {
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Access-Control-Allow-Methods", "*")
	h.Set("Access-Control-Allow-Headers", "*")
}
