package oogway

import (
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
	assert.NoError(t, os.WriteFile(in, []byte(".class{&-name{color:#f00}}"), 0777))
	time.Sleep(time.Millisecond * 10)
	out := filepath.Join(assetsDir, "test.css")
	cfg.Sass.Dir = assetsDir
	cfg.Sass.Entrypoint = "test.scss"
	cfg.Sass.Out = "test.css"
	compileSass()
	assert.FileExists(t, out)
	content, err := os.ReadFile(out)
	assert.NoError(t, err)
	assert.Equal(t, ".class-name{color:red}", string(content))
}
