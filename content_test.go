package oogway

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadContent(t *testing.T) {
	assert.NoError(t, os.RemoveAll(contentDir))
	assert.NoError(t, loadContent(".", nil))
	assert.Len(t, content.templates, 0)
	assert.Len(t, routes.routes, 0)
	assert.NoError(t, os.Mkdir(contentDir, 0777))
	assert.NoError(t, loadContent(".", nil))
	assert.Len(t, content.templates, 0)
	assert.Len(t, routes.routes, 0)
	home := filepath.Join(contentDir, "index.html")
	assert.NoError(t, os.WriteFile(home, []byte("This is the homepage."), 0777))
	assert.NoError(t, loadContent(".", nil))
	assert.Len(t, content.templates, 1)
	assert.Len(t, routes.routes, 1)
	assert.NotNil(t, content.get(home))
	assert.NoError(t, os.Mkdir(filepath.Join(contentDir, "sub"), 0777))
	nested := filepath.Join(contentDir, "sub", "index.html")
	assert.NoError(t, os.WriteFile(nested, []byte("This is a nested page."), 0777))
	assert.NoError(t, loadContent(".", nil))
	assert.Len(t, content.templates, 2)
	assert.Len(t, routes.routes, 2)
	assert.NotNil(t, content.get(nested))
}

func TestWatchContent(t *testing.T) {
	assert.NoError(t, os.RemoveAll(contentDir))
	assert.NoError(t, os.Mkdir(contentDir, 0777))
	home := filepath.Join(contentDir, "index.html")
	assert.NoError(t, os.WriteFile(home, []byte("This is the homepage."), 0777))
	time.Sleep(time.Millisecond * 10)
	ctx, cancel := context.WithCancel(context.Background())
	assert.NoError(t, watchContent(ctx, ".", nil))
	assert.Len(t, content.templates, 1)
	assert.Len(t, routes.routes, 1)
	assert.NoError(t, os.Mkdir(filepath.Join(contentDir, "sub"), 0777))
	nested := filepath.Join(contentDir, "sub", "index.html")
	assert.NoError(t, os.WriteFile(nested, []byte("This is a nested page."), 0777))
	time.Sleep(time.Millisecond * 10)
	assert.Len(t, content.templates, 2)
	assert.Len(t, routes.routes, 2)
	cancel()
}
