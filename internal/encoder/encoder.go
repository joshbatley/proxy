package encoder

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/andybalholm/brotli"
)

// Decompress decodes the body from the compression passed along
func Decompress(h http.Header, compressedBody []byte) ([]byte, error) {
	var content []byte
	var err error
	compressedContent := compressedBody
	br := bytes.NewReader(compressedContent)
	var decompressor io.Reader
	encoding := h.Get("Content-Encoding")

	switch encoding {
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

	content, err = ioutil.ReadAll(decompressor)
	if err != nil {
		return compressedBody, err
	}

	return content, nil
}

// Compress encodes the body from the compression passed along
func Compress(h http.Header, body []byte) ([]byte, error) {
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
