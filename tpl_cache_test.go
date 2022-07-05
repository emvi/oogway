package oogway

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTplCache(t *testing.T) {
	assert.NoError(t, os.RemoveAll(partialsDir))
	assert.NoError(t, os.Mkdir(partialsDir, 0777))
	tplPath := filepath.Join(partialsDir, "tpl.html")
	assert.NoError(t, os.WriteFile(tplPath, []byte("<h1>Hello World!</h1>"), 0777))
	cache := newTplCache()
	assert.Nil(t, cache.get(tplPath))
	addedTpl, err := cache.load(tplPath, nil)
	assert.NoError(t, err)
	assert.NotNil(t, addedTpl)
	tpl := cache.get("tpl")
	assert.NotNil(t, tpl)
	var buffer bytes.Buffer
	assert.NoError(t, tpl.Execute(&buffer, nil))
	assert.Equal(t, "<h1>Hello World!</h1>", buffer.String())
}

func TestTplCacheGetTemplateName(t *testing.T) {
	cache := newTplCache()
	assert.Equal(t, "/", cache.getTemplateName("content/index.html"))
	assert.Equal(t, "/foo", cache.getTemplateName("content/foo/index.html"))
	assert.Equal(t, "test", cache.getTemplateName("partials/test.html"))
	assert.Equal(t, "foo/test", cache.getTemplateName("partials/foo/test.html"))
	assert.Equal(t, "head", cache.getTemplateName("demo/partials/head.html"))
	assert.Equal(t, "", cache.getTemplateName(""))
	assert.Equal(t, "/", cache.getTemplateName("test/content/"))
	assert.Equal(t, "", cache.getTemplateName("test/partials/"))
}
