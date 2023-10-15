package pkg

import (
	"net/http"
	"path/filepath"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
)

const (
	assetsDir = "assets"
)

func serveAssets(router *mux.Router, dir string) {
	fs := http.StripPrefix("/assets/", http.FileServer(http.Dir(filepath.Join(dir, assetsDir))))
	router.PathPrefix("/assets/").Handler(gziphandler.GzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})))
}
