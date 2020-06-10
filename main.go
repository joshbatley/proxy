package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type result struct {
	status  int
	data    []byte
	headers http.Header
}

type request struct {
	method  string
	url     string
	headers http.Header
	body    io.ReadCloser
}

func readBodyToBytes(res io.ReadCloser) []byte {
	body, err := ioutil.ReadAll(res)
	if err != nil {
		panic(err)
	}
	return body
}

func formatURL(url string) string {
	r := strings.Replace(url, "/query?q=", "", 1)
	if strings.Contains(r, "?") {
		return r
	}
	r = strings.Replace(r, "&", "?", 1)
	return r
}

func parseToUnknownJSON(b []byte) map[string]interface{} {
	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		panic(err)
	}
	m := i.(map[string]interface{})
	return m
}

func fetch(r request) result {
	cl := http.Client{}
	req, err := http.NewRequest(r.method, r.url, r.body)
	if err != nil {
		panic(err)
	}
	addResHeaders(r.headers, req.Header)
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	fmt.Println(string(readBodyToBytes(req.Body)))

	res, err := cl.Do(req)
	if err != nil {
		panic(err)
	}

	body := readBodyToBytes(res.Body)

	return result{
		status:  res.StatusCode,
		data:    body,
		headers: res.Header,
	}
}

func addResHeaders(newHeaders http.Header, reqHeaders http.Header) {
	for i, s := range newHeaders {
		if i != "Content-Length" && i != "Cookie" {
			reqHeaders.Set(i, s[0])
		}
	}
}

func query(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-AuthToken, Content-Length, X-Requested-With, X_Auth_Credentials, X-Hub-Version, Cko-Hub-Action")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	toMake := request{
		url:     formatURL(r.URL.String()),
		method:  r.Method,
		headers: r.Header,
		body:    r.Body,
	}

	res := fetch(toMake)
	addResHeaders(res.headers, w.Header())
	w.WriteHeader(res.status)
	w.Write(res.data)
}

func main() {
	http.HandleFunc("/query", query)
	http.ListenAndServe(":8090", nil)
}
