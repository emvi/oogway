package oogway

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCompileSass(t *testing.T) {
	assert.NoError(t, os.RemoveAll(assetsDir))
	assert.NoError(t, os.Mkdir(assetsDir, 0777))
	in := filepath.Join(assetsDir, "test.scss")
	assert.NoError(t, os.WriteFile(in, []byte(".class{&-name{color:red}}"), 0777))
	time.Sleep(time.Millisecond * 10)
	out := filepath.Join(assetsDir, "test.css")
	cfg.Sass.Dir = assetsDir
	cfg.Sass.Entrypoint = "test.scss"
	cfg.Sass.Out = "assets/test.css"
	cfg.Sass.OutSourceMap = "assets/test.css.map"
	compileSass("")
	assert.FileExists(t, out)
	assert.FileExists(t, cfg.Sass.OutSourceMap)
	content, err := os.ReadFile(out)
	assert.NoError(t, err)
	assert.Equal(t, ".class-name{color:red}", string(content))
}

func TestWatchSass(t *testing.T) {
	assert.NoError(t, os.RemoveAll(assetsDir))
	assert.NoError(t, os.Mkdir(assetsDir, 0777))
	in := filepath.Join(assetsDir, "test.scss")
	assert.NoError(t, os.WriteFile(in, []byte(".class{&-name{color:red}}"), 0777))
	time.Sleep(time.Millisecond * 10)
	out := filepath.Join(assetsDir, "test.css")
	cfg.Sass.Dir = assetsDir
	cfg.Sass.Entrypoint = "test.scss"
	cfg.Sass.Out = "assets/test.css"
	cfg.Sass.Watch = true
	ctx, cancel := context.WithCancel(context.Background())
	assert.NoError(t, watchSass(ctx, ""))
	time.Sleep(time.Millisecond * 10)
	assert.FileExists(t, out)
	content, err := os.ReadFile(out)
	assert.NoError(t, err)
	assert.Equal(t, ".class-name{color:red}", string(content))
	assert.NoError(t, os.WriteFile(in, []byte(".class{&-name{color:blue}}"), 0777))
	time.Sleep(time.Millisecond * 10)
	assert.FileExists(t, out)
	content, err = os.ReadFile(out)
	assert.NoError(t, err)
	assert.Equal(t, ".class-name{color:blue}", string(content))
	cancel()
}
