package handler

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ClientServe - Serve a webapp
func ClientServe(w http.ResponseWriter, r *http.Request) {
	StaticPath := "./webapp/build"
	IndexPath := "index.html"

	// get the absolute path to prevent directory traversal
	path := strings.Replace(r.URL.Path, "/config", "", 1)
	path, err := filepath.Abs(path)

	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(StaticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		log.Println("index")
		http.ServeFile(w, r, filepath.Join(StaticPath, IndexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.StripPrefix(
		"/config/",
		http.FileServer(http.Dir(StaticPath)),
	).ServeHTTP(w, r)
}
