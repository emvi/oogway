package oogway

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
)

const (
	contentDir      = "content"
	contentPageFile = "index.html"
)

var (
	tpl    = newTplCache()
	routes = newRouter()
)

func loadContent(dir string) error {
	contentDirPath := filepath.Join(dir, contentDir)

	if _, err := os.Stat(contentDirPath); os.IsNotExist(err) || isEmptyDir(contentDirPath) {
		return nil
	}

	return filepath.WalkDir(contentDirPath, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && d.Name() == contentPageFile {
			tpl, err := tpl.load(path)

			if err != nil {
				log.Printf("Error loading template %s: %s", path, err)
				return nil
			}

			route := filepath.Dir(path)[len(contentDirPath):] + "/"
			routes.addRoute(route, tpl)
			log.Printf("Added route %s", route)
		}

		return nil
	})
}

func servePage(router *mux.Router) {
	router.PathPrefix("/").Handler(gziphandler.GzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if !strings.HasSuffix(path, "/") {
			path += "/"
		}

		tpl := routes.findTemplate(path)

		if tpl == nil {
			// TODO configurable 404 page
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err := tpl.Execute(w, nil); err != nil {
			log.Printf("Error rendering page %s: %s", r.URL.Path, err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})))
}
