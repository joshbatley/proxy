package client

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Handler has the correct rules a SPA
type Handler struct {
	StaticPath string
	IndexPath  string
}

// ServeHTTP sets up SPA endpoint
func (c Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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
	path = filepath.Join(c.StaticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(c.StaticPath, c.IndexPath))
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
		http.FileServer(http.Dir(c.StaticPath)),
	).ServeHTTP(w, r)
}
