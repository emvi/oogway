package oogway

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadContent(t *testing.T) {
	tpl.clear()
	assert.NoError(t, os.RemoveAll(contentDir))
	assert.NoError(t, loadContent("."))
	assert.Len(t, tpl.templates, 0)
	assert.NoError(t, os.Mkdir(contentDir, 0777))
	assert.NoError(t, loadContent("."))
	assert.Len(t, tpl.templates, 0)
	home := filepath.Join(contentDir, "index.html")
	assert.NoError(t, os.WriteFile(home, []byte("This is the homepage."), 0777))
	assert.NoError(t, loadContent("."))
	assert.Len(t, tpl.templates, 1)
	assert.NotNil(t, tpl.get(home))
	assert.NoError(t, os.Mkdir(filepath.Join(contentDir, "sub"), 0777))
	nested := filepath.Join(contentDir, "sub", "index.html")
	assert.NoError(t, os.WriteFile(nested, []byte("This is a nested page."), 0777))
	assert.NoError(t, loadContent("."))
	assert.Len(t, tpl.templates, 2)
	assert.NotNil(t, tpl.get(nested))
}
