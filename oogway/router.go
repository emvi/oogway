package oogway

import (
	"html/template"
	"sync"
)

type router struct {
	routes map[string]*template.Template
	m      sync.RWMutex
}

func newRouter() *router {
	return &router{
		routes: make(map[string]*template.Template),
	}
}

func (router *router) addRoute(route string, template *template.Template) {
	router.m.Lock()
	defer router.m.Unlock()
	router.routes[route] = template
}

func (router *router) findTemplate(route string) *template.Template {
	router.m.RLock()
	defer router.m.RUnlock()
	return router.routes[route]
}

func (router *router) clear() {
	router.m.Lock()
	defer router.m.Unlock()
	router.routes = make(map[string]*template.Template)
}
