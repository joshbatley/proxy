package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type config struct {
	Name string `yml:"name"`
	Port string `yml:"port"`
}

type req struct {
	url     *url.URL
	body    []byte
	headers http.Header
	status  int
	method  string
}

var cache []req

func readBodyToBytes(res io.ReadCloser) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(res)
	return buf.Bytes()
}

func formatURL(dirtyURL string) *url.URL {
	s := regexp.MustCompile(`(?:/query\?q=)(.{0,})`)
	r := string(s.ReplaceAll([]byte(dirtyURL), []byte("$1")))
	formattedURL, err := url.Parse(r)
	if err != nil {
		panic(err)
	}
	return formattedURL
}

func parseToUnknownJSON(b []byte) map[string]interface{} {
	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		panic(err)
	}
	m := i.(map[string]interface{})
	return m
}

func getPreRequest(url *url.URL, r *http.Request) bool {
	fmt.Println(url, r.Method)
	matched, err := regexp.Match("posts", []byte(url.String()))
	if err != nil {
		panic(err)
	}
	if matched {
		return true
	}
	return false
}

func findInCache(url *url.URL) (req, bool) {
	for _, c := range cache {
		if strings.Compare(c.url.String(), url.String()) >= 0 {
			return c, true
		}
	}
	return req{}, false
}

func postRequest(res *http.Response) error {
	fmt.Println("caching")

	buf, _ := ioutil.ReadAll(res.Body)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

	newC := req{
		url:     res.Request.URL,
		body:    readBodyToBytes(rdr1),
		headers: res.Header,
		method:  res.Request.Method,
		status:  res.StatusCode,
	}
	cache = append(cache, newC)
	res.Body = rdr2
	fmt.Println(res.Status, res.StatusCode)
	return nil
}

func query(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, PUT, DELETE, CONNECT, OPTIONS, TRACE, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "*, Origin, X-Requested-With, Content-Type, Accept")
		return
	}

	url := formatURL(r.URL.String())

	// if returnEarly := getPreRequest(url, r); returnEarly {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	fmt.Println(url)
	data, found := findInCache(url)

	if found == true {
		fmt.Println("found in cache sending cache")
		for i, h := range data.headers {
			w.Header().Set(i, strings.Join(h, " "))
		}
		w.WriteHeader(data.status)
		w.Write(data.body)
		return
	}

	fmt.Println("no entry fetching from server")

	reverseProxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = url.Scheme
			req.URL.Host = url.Host
			req.URL.Path = url.Path
			req.Host = url.Host
			req.URL.RawQuery = url.RawQuery
		},
		ModifyResponse: postRequest,
	}

	reverseProxy.ServeHTTP(w, r)

}

func readConfig() (*config, error) {
	file, _ := ioutil.ReadFile("./config.yml")
	config := &config{}
	err := yaml.Unmarshal([]byte(file), &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func main() {
	config, err := readConfig()
	if err != nil {
		panic("Config unreadable")
	}
	http.HandleFunc("/query", query)
	http.ListenAndServe(":"+config.Port, nil)

	fmt.Println("listing on 127.0.0.1:" + config.Port)
}
