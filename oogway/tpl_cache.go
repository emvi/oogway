package oogway

import (
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var (
	removeBasePath = regexp.MustCompile("^(.*(?:partials|content))/")
)

type tplCache struct {
	templates map[string]template.Template
	m         sync.RWMutex
}

func newTplCache() *tplCache {
	return &tplCache{
		templates: make(map[string]template.Template),
	}
}

func (cache *tplCache) load(path string, funcMap template.FuncMap) (*template.Template, error) {
	cache.m.Lock()
	defer cache.m.Unlock()
	content, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	tpl, err := template.New(path).Funcs(funcMap).Parse(string(content))

	if err != nil {
		return nil, err
	}

	cache.templates[cache.getTemplateName(path)] = *tpl
	return tpl, nil
}

func (cache *tplCache) get(name string) *template.Template {
	cache.m.RLock()
	defer cache.m.RUnlock()
	tpl, found := cache.templates[name]

	if !found {
		return nil
	}

	return &tpl
}

func (cache *tplCache) clear() {
	cache.m.Lock()
	defer cache.m.Unlock()
	cache.templates = make(map[string]template.Template)
}

func (cache *tplCache) getTemplateName(path string) string {
	path = strings.ReplaceAll(path, "\\", "/")
	base := removeBasePath.FindStringSubmatch(path)

	if len(base) != 2 {
		return ""
	}

	path = path[len(base[0]):]

	if strings.HasSuffix(base[1], contentDir) {
		path = filepath.Dir(path)

		if path == "" || path == "." {
			path = "/"
		}

		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
	} else if strings.HasSuffix(base[1], partialsDir) {
		path = strings.TrimSuffix(path, ".html")
	}

	return path
}
