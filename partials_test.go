package oogway

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadPartials(t *testing.T) {
	assert.NoError(t, os.RemoveAll(partialsDir))
	assert.NoError(t, loadPartials(".", nil))
	assert.NoError(t, os.Mkdir(partialsDir, 0777))
	assert.NoError(t, loadPartials(".", nil))
	tplPath := filepath.Join(partialsDir, "tpl.html")
	assert.NoError(t, os.WriteFile(tplPath, []byte(`<h1>Hello World!</h1>`), 0777))
	time.Sleep(time.Millisecond * 10)
	assert.NoError(t, loadPartials(".", nil))
	var buffer bytes.Buffer
	assert.NoError(t, partials.get(tplPath).Execute(&buffer, nil))
	assert.Equal(t, "<h1>Hello World!</h1>", buffer.String())
}

func TestWatchPartials(t *testing.T) {
	assert.NoError(t, os.RemoveAll(partialsDir))
	assert.NoError(t, os.Mkdir(partialsDir, 0777))
	tplPath := filepath.Join(partialsDir, "tpl.html")
	assert.NoError(t, os.WriteFile(tplPath, []byte(`<h1>Hello World!</h1>`), 0777))
	time.Sleep(time.Millisecond * 10)
	ctx, cancel := context.WithCancel(context.Background())
	assert.NoError(t, watchPartials(ctx, ".", nil))
	var buffer bytes.Buffer
	assert.NoError(t, partials.get(tplPath).Execute(&buffer, nil))
	assert.Equal(t, "<h1>Hello World!</h1>", buffer.String())
	assert.NoError(t, os.WriteFile(tplPath, []byte(`<p>Lorem ipsum</p>`), 0777))
	time.Sleep(time.Millisecond * 10)
	buffer.Reset()
	assert.NoError(t, partials.get(tplPath).Execute(&buffer, nil))
	assert.Equal(t, "<p>Lorem ipsum</p>", buffer.String())
	cancel()
}
