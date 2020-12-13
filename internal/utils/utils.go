package utils

import (
	"encoding/json"
	"net/http"
	"reflect"

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

type response struct {
	Count int         `json:"count"`
	Skip  int         `json:"skip"`
	Limit int         `json:"limit"`
	Data  interface{} `json:"data"`
}

// Response -
func Response(w http.ResponseWriter, d interface{}, l int, s int) {
	var count int
	if reflect.TypeOf(d).Kind() == reflect.Slice {
		count = reflect.ValueOf(d).Len()
	}
	res := response{
		Count: count,
		Skip:  s,
		Limit: l,
		Data:  d,
	}

	j, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json, text/plain, */*")
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}
