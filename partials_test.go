package oogway

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadPartials(t *testing.T) {
	tpl.clear()
	assert.NoError(t, os.RemoveAll(partialsDir))
	assert.NoError(t, loadPartials("."))
	assert.NoError(t, os.Mkdir(partialsDir, 0777))
	assert.NoError(t, loadPartials("."))
	tplPath := filepath.Join(partialsDir, "tpl.html")
	assert.NoError(t, os.WriteFile(tplPath, []byte(`<h1>Hello World!</h1>`), 0777))
	time.Sleep(time.Millisecond * 10)
	assert.NoError(t, loadPartials("."))
	var buffer bytes.Buffer
	assert.NoError(t, tpl.get(tplPath).Execute(&buffer, nil))
	assert.Equal(t, "<h1>Hello World!</h1>", buffer.String())
}
