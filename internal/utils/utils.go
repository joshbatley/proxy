package utils

import (
	"encoding/json"
	"net/http"

	"github.com/joshbatley/proxy/internal/fail"
)

// NotFound 404 handler
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json, text/plain, */*")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{}`))
}

// BadRequest 400 response
func BadRequest(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json, text/plain, */*")
	w.WriteHeader(http.StatusBadRequest)
	jsonString, _ := json.Marshal(err)

	if len(jsonString) == 2 {
		jsonString, _ = json.Marshal(fail.InternalError(err))
	}

	w.Write(jsonString)
}

// Cors -
func Cors(h http.Header) {
	h.Set("Access-Control-Allow-Origin", "*")
	h.Set("Access-Control-Allow-Methods", "*")
	h.Set("Access-Control-Allow-Headers", "*")
}
