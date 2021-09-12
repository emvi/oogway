package oogway

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
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
	assert.NotNil(t, content.get("/"))
	assert.NoError(t, os.Mkdir(filepath.Join(contentDir, "sub"), 0777))
	nested := filepath.Join(contentDir, "sub", "index.html")
	assert.NoError(t, os.WriteFile(nested, []byte("This is a nested page."), 0777))
	assert.NoError(t, loadContent(".", nil))
	assert.Len(t, content.templates, 2)
	assert.Len(t, routes.routes, 2)
	assert.NotNil(t, content.get("/sub"))
}

func TestWatchContent(t *testing.T) {
	assert.NoError(t, os.RemoveAll(contentDir))
	assert.NoError(t, os.Mkdir(contentDir, 0777))
	home := filepath.Join(contentDir, "index.html")
	assert.NoError(t, os.WriteFile(home, []byte("This is the homepage."), 0777))
	time.Sleep(time.Millisecond * 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	assert.NoError(t, watchContent(ctx, ".", nil))
	assert.Len(t, content.templates, 1)
	assert.Len(t, routes.routes, 1)
	assert.NoError(t, os.Mkdir(filepath.Join(contentDir, "sub"), 0777))
	nested := filepath.Join(contentDir, "sub", "index.html")
	assert.NoError(t, os.WriteFile(nested, []byte("This is a nested page."), 0777))
	time.Sleep(time.Millisecond * 10)
	assert.Len(t, content.templates, 2)
	assert.Len(t, routes.routes, 2)
}

func TestRenderPage(t *testing.T) {
	assert.NoError(t, os.RemoveAll(contentDir))
	assert.NoError(t, os.MkdirAll(filepath.Join(contentDir, "not-found"), 0777))
	home := filepath.Join(contentDir, "index.html")
	assert.NoError(t, os.WriteFile(home, []byte("This is the homepage."), 0777))
	notFound := filepath.Join(contentDir, "not-found", "index.html")
	assert.NoError(t, os.WriteFile(notFound, []byte("Page not found."), 0777))
	time.Sleep(time.Millisecond * 10)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	assert.NoError(t, watchContent(ctx, ".", nil))

	// request homepage
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	renderPage(w, r)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	body, _ := io.ReadAll(w.Result().Body)
	assert.Equal(t, "This is the homepage.", string(body))

	// page not found (no 404 page)
	r = httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	w = httptest.NewRecorder()
	renderPage(w, r)
	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	body, _ = io.ReadAll(w.Result().Body)
	assert.Empty(t, string(body))

	// page not found (404 page configured)
	cfg.Content.NotFound = "/not-found"
	r = httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	w = httptest.NewRecorder()
	renderPage(w, r)
	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	body, _ = io.ReadAll(w.Result().Body)
	assert.Equal(t, "Page not found.", string(body))
}
