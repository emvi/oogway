package oogway

import (
	"context"
	"github.com/stretchr/testify/assert"
	"html/template"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestMergeFuncMaps(t *testing.T) {
	a := template.FuncMap{
		"a": func() int { return 42 },
	}
	b := template.FuncMap{
		"a": func() int { return 43 },
		"b": func() int { return 44 },
	}
	out := mergeFuncMaps(a, b)
	assert.NotNil(t, out["a"])
	assert.NotNil(t, out["b"])
	assert.Nil(t, out["c"])
}

func TestRenderTemplate(t *testing.T) {
	assert.NoError(t, os.RemoveAll(contentDir))
	assert.NoError(t, os.Mkdir(contentDir, 0777))
	home := filepath.Join(contentDir, "index.html")
	assert.NoError(t, os.WriteFile(home, []byte("This is the homepage."), 0777))
	time.Sleep(time.Millisecond * 10)
	ctx, cancel := context.WithCancel(context.Background())
	assert.NoError(t, watchContent(ctx, ".", nil))
	assert.Equal(t, template.HTML("This is the homepage."), renderTemplate("content/index.html", nil))
	cancel()
}

func TestRenderPartial(t *testing.T) {
	assert.NoError(t, os.RemoveAll(partialsDir))
	assert.NoError(t, os.Mkdir(partialsDir, 0777))
	partial := filepath.Join(partialsDir, "partial.html")
	assert.NoError(t, os.WriteFile(partial, []byte("This is a partial."), 0777))
	time.Sleep(time.Millisecond * 10)
	ctx, cancel := context.WithCancel(context.Background())
	assert.NoError(t, watchPartials(ctx, ".", nil))
	assert.Equal(t, template.HTML("This is a partial."), renderPartial("partials/partial.html", nil))
	cancel()
}
