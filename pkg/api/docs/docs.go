package docs

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

// DocsDIR constant represents where documentation files are stored
const DocsDIR = "resources/docs"

// DocsPATH constant represents http path where documentation files can be reached
const DocsPATH = "/docs"

// DocServer serves documentation files
func DocServer(r chi.Router, path string, root http.FileSystem) {
	serveDocs(r, path, root)
}

// serveDocs conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func serveDocs(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
