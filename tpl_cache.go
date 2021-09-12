package oogway

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
	if strings.HasPrefix(path, contentDir) {
		path = filepath.Dir(path)
		path = path[len(contentDir):]

		if path == "" {
			path = "/"
		}
	} else if strings.HasPrefix(path, partialsDir+"/") {
		path = path[len(partialsDir+"/"):]

		if strings.HasSuffix(path, ".html") {
			path = path[:len(path)-len(".html")]
		}
	}

	return path
}
