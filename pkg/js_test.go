package pkg

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCompileJS(t *testing.T) {
	assert.NoError(t, os.RemoveAll(assetsDir))
	assert.NoError(t, os.Mkdir(assetsDir, 0777))
	in := filepath.Join(assetsDir, "test.js")
	assert.NoError(t, os.WriteFile(in, []byte("console.log('Hello World')"), 0777))
	time.Sleep(time.Millisecond * 10)
	out := filepath.Join(assetsDir, "bundle.js")
	cfg.JS.Dir = assetsDir
	cfg.JS.Entrypoint = "test.js"
	cfg.JS.Out = "assets/bundle.js"
	cfg.JS.SourceMap = true
	compileJS("")
	assert.FileExists(t, out)
	assert.FileExists(t, "assets/bundle.js.map")
	content, err := os.ReadFile(out)
	assert.NoError(t, err)
	assert.Equal(t, "(()=>{console.log(\"Hello World\");})();\n", string(content))
}

func TestWatchJS(t *testing.T) {
	assert.NoError(t, os.RemoveAll(assetsDir))
	assert.NoError(t, os.Mkdir(assetsDir, 0777))
	in := filepath.Join(assetsDir, "test.js")
	assert.NoError(t, os.WriteFile(in, []byte("console.log('Hello World')"), 0777))
	time.Sleep(time.Millisecond * 10)
	out := filepath.Join(assetsDir, "bundle.js")
	cfg.JS.Dir = assetsDir
	cfg.JS.Entrypoint = "test.js"
	cfg.JS.Out = "assets/bundle.js"
	cfg.JS.Watch = true
	ctx, cancel := context.WithCancel(context.Background())
	assert.NoError(t, watchJS(ctx, ""))
	time.Sleep(time.Millisecond * 10)
	assert.FileExists(t, out)
	content, err := os.ReadFile(out)
	assert.NoError(t, err)
	assert.Equal(t, "(()=>{console.log(\"Hello World\");})();\n", string(content))
	assert.NoError(t, os.WriteFile(in, []byte("console.log('Foo bar')"), 0777))
	time.Sleep(time.Millisecond * 10)
	assert.FileExists(t, out)
	content, err = os.ReadFile(out)
	assert.NoError(t, err)
	assert.Equal(t, "(()=>{console.log(\"Foo bar\");})();\n", string(content))
	cancel()
}
