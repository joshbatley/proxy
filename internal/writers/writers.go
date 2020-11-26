package writers

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/andybalholm/brotli"
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

// DecodeBody decodes the body from the compression passed along
func DecodeBody(h http.Header, compressedBody []byte) ([]byte, error) {
	var content []byte
	var err error
	compressedContent := compressedBody
	br := bytes.NewReader(compressedContent)
	var decompressor io.Reader

	if ce := h.Get("Content-Encoding"); ce != "" {
		switch ce {
		case "br":
			decompressor = brotli.NewReader(br)
		case "deflate":
			decompressor = flate.NewReader(br)
		case "gzip":
			decompressor, err = gzip.NewReader(br)
			if err != nil {
				decompressor = br
			}
		default:
			decompressor = br
		}
	}

	content, err = ioutil.ReadAll(decompressor)
	if err != nil {
		return compressedBody, err
	}

	return content, nil
}

// EncodeBody encodes the body from the compression passed along
func EncodeBody(h http.Header, body []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	var compressor io.WriteCloser
	var compressedContent []byte
	var err error
	encoding := h.Get("Content-Encoding")

	switch encoding {
	case "br":
		compressor = brotli.NewWriter(buf)
	case "gzip":
		compressor = gzip.NewWriter(buf)
		if err != nil {
			log.Println("Error creating gzip compressor:", err)
			compressor = nil
		}
	default:
		compressedContent = body
	}
	if compressor != nil {
		compressor.Write(body)
		if err := compressor.Close(); err != nil {
			log.Println("Error compressing content of")
		} else {
			compressedContent = buf.Bytes()
		}
	}

	return compressedContent, nil
}
