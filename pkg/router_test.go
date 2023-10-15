package pkg

import (
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	router := newRouter()
	router.addRoute("/", template.New("home"))
	router.addRoute("/foo", template.New("foo"))
	router.addRoute("/foo/bar", template.New("bar"))
	assert.Nil(t, router.findTemplate("/not-found"))
	assert.NotNil(t, router.findTemplate("/"))
	assert.NotNil(t, router.findTemplate("/foo"))
	assert.NotNil(t, router.findTemplate("/foo/bar"))
	router.clear()
	assert.Nil(t, router.findTemplate("/"))
	assert.Nil(t, router.findTemplate("/foo"))
	assert.Nil(t, router.findTemplate("/foo/bar"))
}
