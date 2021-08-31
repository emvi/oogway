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
	assert.NoError(t, os.RemoveAll(partialsDir))
	assert.NoError(t, loadPartials("."))
	assert.NoError(t, os.Mkdir(partialsDir, 0777))
	assert.NoError(t, loadPartials("."))
	assert.NoError(t, os.WriteFile(filepath.Join(partialsDir, "tpl.html"), []byte(`<h1>Hello World!</h1>`), 0777))
	time.Sleep(time.Millisecond * 10)
	assert.NoError(t, loadPartials("."))
	var buffer bytes.Buffer
	assert.NoError(t, tpl.ExecuteTemplate(&buffer, "tpl.html", nil))
	assert.Equal(t, "<h1>Hello World!</h1>", buffer.String())
}
