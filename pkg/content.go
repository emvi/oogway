package pkg

import (
	"context"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/fsnotify/fsnotify"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
)

const (
	contentDir      = "content"
	contentPageFile = "index.html"
	metaPageFile    = "meta.toml"
)

type meta struct {
	SitemapPriority float64 `toml:"sitemap_priority"`
}

var (
	content    = newTplCache()
	routes     = newRouter()
	sitemap    []sitemapURL
	sitemapXML []byte
)

func loadContent(dir string, funcMap template.FuncMap) error {
	content.clear()
	routes.clear()
	sitemap = make([]sitemapURL, 0)
	contentDirPath := filepath.Join(dir, contentDir)

	if _, err := os.Stat(contentDirPath); os.IsNotExist(err) || isEmptyDir(contentDirPath) {
		return nil
	}

	if err := filepath.WalkDir(contentDirPath, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && d.Name() == contentPageFile {
			tpl, err := content.load(path, funcMap)

			if err != nil {
				log.Printf("Error loading template %s: %s", path, err)
				return nil
			}

			route := filepath.Dir(path)[len(contentDirPath):] + "/"
			route = strings.ReplaceAll(route, "\\", "/")
			routes.addRoute(route, tpl)
			m := loadMetaInformation(filepath.Dir(path))
			sitemapPriority := m.SitemapPriority

			if sitemapPriority <= 0.001 {
				sitemapPriority = 1
			}

			sitemap = append(sitemap, sitemapURL{
				Loc:      route,
				Priority: strconv.FormatFloat(sitemapPriority, 'f', 2, 64),
				Lastmod:  time.Now().Format(sitemapLastModFormat),
			})
		}

		return nil
	}); err != nil {
		return err
	}

	var err error
	sitemapXML, err = generateSitemap(sitemap)
	return err
}

func watchContent(ctx context.Context, dir string, funcMap template.FuncMap) error {
	if err := loadContent(dir, funcMap); err != nil {
		return err
	}

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case _, ok := <-watcher.Events:
				if !ok {
					continue
				}

				if err := loadContent(dir, funcMap); err != nil {
					log.Printf("Error updating content: %s", err)
				}
			case <-ctx.Done():
				watcher.Close()
				return
			}
		}
	}()

	if err := watcher.Add(filepath.Join(dir, contentDir)); err != nil {
		return err
	}

	return nil
}

func loadMetaInformation(path string) meta {
	var m meta
	content, err := os.ReadFile(filepath.Join(path, metaPageFile))

	if os.IsNotExist(err) {
		return m
	}

	if err != nil {
		log.Printf("Error opening meta file %s: %s", path, err)
		return m
	}

	if err := toml.Unmarshal(content, &m); err != nil {
		log.Printf("Error loading meta file %s: %s", path, err)
		return m
	}

	return m
}

func servePage(router *mux.Router) {
	router.PathPrefix("/").Handler(gziphandler.GzipHandler(http.HandlerFunc(renderPage)))
}

func serveSitemap(router *mux.Router) {
	router.Path("/sitemap.xml").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write(sitemapXML); err != nil {
			log.Printf("Error serving sitemap: %s", err)
		}
	})
}

func renderPage(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	tpl := routes.findTemplate(path)

	if tpl == nil {
		w.WriteHeader(http.StatusNotFound)

		if cfg.Content.NotFound != "" {
			path = cfg.Content.NotFound

			if !strings.HasSuffix(path, "/") {
				path += "/"
			}

			tpl = routes.findTemplate(path)

			if tpl != nil {
				go pageView(r, cfg.Content.NotFound)

				if err := tpl.Execute(w, nil); err != nil {
					log.Printf("Error rendering page %s: %s", r.URL.Path, err)
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}

		return
	}

	go pageView(r, "")

	if err := tpl.Execute(w, nil); err != nil {
		log.Printf("Error rendering page %s: %s", r.URL.Path, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
