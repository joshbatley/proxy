package query

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// Req - Request Data struct
type Req struct {
	URL     *url.URL
	Body    []byte
	Headers http.Header
	Status  int
	Method  string
}

// GetPreResponse -
func getPreResponse(url *url.URL, r *http.Request, w http.ResponseWriter) bool {
	log.Println(url, r.Method)
	matched, err := regexp.Match("posts", []byte(url.String()))
	if err != nil {
		panic(err)
	}
	if matched {
		b := []byte{}
		w.WriteHeader(http.StatusNoContent)
		w.Write(b)
		return true
	}
	return false
}

// ModifyResponse -
func modifyResponse(cache []Req) func(res *http.Response) error {
	return func(res *http.Response) error {
		log.Println("caching")

		buf, _ := ioutil.ReadAll(res.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

		newC := Req{
			URL:     res.Request.URL,
			Body:    readBodyToBytes(rdr1),
			Headers: res.Header,
			Method:  res.Request.Method,
			Status:  res.StatusCode,
		}
		cache = append(cache, newC)
		res.Body = rdr2
		log.Println(res.Status, res.StatusCode)
		return nil
	}
}

// HandleOptions -
func handleOptions(w http.ResponseWriter, m string) bool {
	if m == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, CONNECT, OPTIONS, TRACE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "*, Origin, X-Requested-With, Content-Type, Accept")
		w.WriteHeader(http.StatusNoContent)
		return true
	}
	return false
}

// SendCache -
func sendCache(url *url.URL, w http.ResponseWriter, cache []Req) bool {
	if data, found := findInCache(url.String(), cache); found == true {
		log.Println("found in cache sending cache")
		for i, h := range data.Headers {
			w.Header().Set(i, strings.Join(h, " "))
		}
		w.WriteHeader(data.Status)
		w.Write(data.Body)
		return true
	}
	return false
}

//
//
func readBodyToBytes(res io.ReadCloser) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(res)
	return buf.Bytes()
}

func findInCache(url string, arr []Req) (Req, bool) {
	for _, c := range arr {
		if strings.Compare(c.URL.String(), url) >= 0 {
			return c, true
		}
	}
	return Req{}, false
}
