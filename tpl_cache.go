package oogway

import (
	"html/template"
	"os"
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

func (cache *tplCache) load(path string) (*template.Template, error) {
	cache.m.Lock()
	defer cache.m.Unlock()
	content, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	tpl, err := template.New(path).Parse(string(content))

	if err != nil {
		return nil, err
	}

	cache.templates[path] = *tpl
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
