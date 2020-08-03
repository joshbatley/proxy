package service

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	database "github.com/joshbatley/proxy/database"
	"github.com/joshbatley/proxy/def"
)

// FormatURL -
func FormatURL(u string) *url.URL {
	s := regexp.MustCompile(`(?:/query\?q=)(.{0,})`)
	r := string(s.ReplaceAll([]byte(u), []byte("$1")))

	formattedURL, err := url.Parse(r)

	if err != nil {
		panic(err)
	}

	return formattedURL
}

type res struct {
	URL     string
	Body    string
	Headers string
	Status  int
	Method  string
}

// func sendCache2(url *url.URL, w http.ResponseWriter) bool {
// 	var data res
// 	log.Print(url.String())
// 	log.Print(url.String())

// 	row := db.QueryRow("SELECT url, body, headers, status, method FROM cache WHERE url = '?'", url.String())
// 	err := row.Scan(&data.URL, &data.Body, &data.Headers, &data.Status, &data.Method)
// 	if err == (sql.ErrNoRows) {
// 		log.Println("no row")
// 		return false
// 	}
// 	if err != nil {
// 		log.Println(err)
// 		return false
// 	}

// 	if data != (res{}) {
// 		log.Println("found in cache sending cache")

// 		for _, h := range strings.Split(data.Headers, ";") {
// 			v := strings.Split(h, "=[")
// 			w.Header().Set(v[0], strings.Replace(v[1], "]", "", 1))
// 		}

// 		w.WriteHeader(data.Status)
// 		w.Write([]byte(data.Body))

// 		return true

// 	}
// 	return false
// }

func checkCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			log.Fatal(err)
		}
	}
	return count
}

// GetPreResponse -
func GetPreResponse(url *url.URL, r *http.Request, w http.ResponseWriter) bool {
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

// // CreateCache =
// func CreateCache(db *sql.DB, r Record) {
// 	tx, _ := db.Begin()
// 	stmt, _ := tx.Prepare(`INSERT INTO cache (collection, url, headers, body, status, method) values (?,?,?,?,?,?)`)

// 	b := new(bytes.Buffer)
// 	for key, value := range r.Headers {
// 		fmt.Fprintf(b, "%s=%s;\n", key, value)
// 	}

// 	_, err := stmt.Exec(0, r.URL.String(), b.String(), r.Body, r.Status, r.Method)
// 	if err != nil {
// 		panic(err)
// 	}
// 	tx.Commit()
// }

// ModifyResponse -
func ModifyResponse(res *http.Response) error {

	log.Println("caching")

	buf, _ := ioutil.ReadAll(res.Body)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

	body := readBodyToBytes(rdr1)

	b := new(bytes.Buffer)
	for key, value := range res.Header {
		fmt.Fprintf(b, "%s=%s;\n", key, value)
	}

	newC := def.Record{
		URL:     res.Request.URL,
		Body:    body,
		Headers: res.Header,
		Method:  res.Request.Method,
		Status:  res.StatusCode,
	}
	database.Insert(newC)
	res.Body = rdr2
	log.Println(res.Status, res.StatusCode)
	return nil

}

// SendCache -
func SendCache(url *url.URL, w http.ResponseWriter, cache []def.Record) bool {
	if data, found := findInCache(url.String(), cache); found == true {
		log.Println("found in cache sending cache")

		// for i, h := range data.Headers {
		// 	w.Header().Set(i, strings.Join(h, " "))
		// }

		w.WriteHeader(data.Status)
		w.Write(data.Body)

		return true

	}

	return false
}

func readBodyToBytes(res io.ReadCloser) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(res)
	return buf.Bytes()
}

func findInCache(url string, arr []def.Record) (def.Record, bool) {
	for _, c := range arr {
		if strings.Compare(c.URLString(), url) >= 0 {
			return c, true
		}
	}
	return def.Record{}, false
}
