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
	assert.NoError(t, cache.load(tplPath))
	tpl := cache.get(tplPath)
	assert.NotNil(t, tpl)
	var buffer bytes.Buffer
	assert.NoError(t, tpl.Execute(&buffer, nil))
	assert.Equal(t, "<h1>Hello World!</h1>", buffer.String())
}
