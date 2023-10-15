package pkg

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	sampleConfig = `[server]
host = "localhost"
port = 8080
shutdown_time = 9
write_timeout = 10
read_timeout = 11
tls_cert_file = "cert/file.pem"
tls_key_file = "key/file.pem"

[content]
not_found = "not-found"

[sass]
entrypoint = "style.scss"
dir = "assets/scss"
watch = true
out = "assets/style.css"
out_source_map = "assets/style.css.map"

[js]
entrypoint = "main.js"
dir = "assets/js"
watch = true
out = "assets/bundle.js"
source_map = true

[pirsch]
client_id = "id"
client_secret = "secret"
subnets = ["10.1.0.0/16", "10.2.0.0/8"]
header = ["X-Forwarded-For", "Forwarded"]
`
)

func TestLoadConfig(t *testing.T) {
	assert.NoError(t, os.RemoveAll("config.toml"))
	assert.NoError(t, os.WriteFile("config.toml", []byte(sampleConfig), 0644))
	assert.NoError(t, loadConfig("."))
	assert.Equal(t, "localhost", cfg.Server.Host)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, 9, cfg.Server.ShutdownTimeout)
	assert.Equal(t, 10, cfg.Server.WriteTimeout)
	assert.Equal(t, 11, cfg.Server.ReadTimeout)
	assert.Equal(t, "cert/file.pem", cfg.Server.TLSCertFile)
	assert.Equal(t, "key/file.pem", cfg.Server.TLSKeyFile)
	assert.Equal(t, "not-found", cfg.Content.NotFound)
	assert.Equal(t, "style.scss", cfg.Sass.Entrypoint)
	assert.Equal(t, "assets/scss", cfg.Sass.Dir)
	assert.True(t, cfg.Sass.Watch)
	assert.Equal(t, "assets/style.css", cfg.Sass.Out)
	assert.Equal(t, "assets/style.css.map", cfg.Sass.OutSourceMap)
	assert.Equal(t, "main.js", cfg.JS.Entrypoint)
	assert.Equal(t, "assets/js", cfg.JS.Dir)
	assert.True(t, cfg.JS.Watch)
	assert.Equal(t, "assets/bundle.js", cfg.JS.Out)
	assert.True(t, cfg.JS.SourceMap)
	assert.Equal(t, "id", cfg.Pirsch.ClientID)
	assert.Equal(t, "secret", cfg.Pirsch.ClientSecret)
	assert.Len(t, cfg.Pirsch.Subnets, 2)
	assert.Equal(t, "10.1.0.0/16", cfg.Pirsch.Subnets[0])
	assert.Equal(t, "10.2.0.0/8", cfg.Pirsch.Subnets[1])
	assert.Len(t, cfg.Pirsch.Header, 2)
	assert.Equal(t, "X-Forwarded-For", cfg.Pirsch.Header[0])
	assert.Equal(t, "Forwarded", cfg.Pirsch.Header[1])
}

func TestLoadConfigNotExists(t *testing.T) {
	assert.NoError(t, os.RemoveAll("config.toml"))
	err := loadConfig(".")
	assert.NotNil(t, err)
	assert.Equal(t, "error loading config.toml: open config.toml: no such file or directory", err.Error())
}

func TestWatchConfig(t *testing.T) {
	assert.NoError(t, os.RemoveAll("config.toml"))
	assert.NoError(t, os.WriteFile("config.toml", []byte(sampleConfig), 0644))
	ctx, cancel := context.WithCancel(context.Background())
	assert.NoError(t, watchConfig(ctx, "."))
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.NoError(t, os.WriteFile("config.toml", []byte(strings.Replace(sampleConfig, "8080", "8888", 1)), 0644))
	time.Sleep(time.Millisecond * 10)
	cancel()
	assert.Equal(t, 8888, cfg.Server.Port)
}

func TestWatchConfigNotExists(t *testing.T) {
	assert.NoError(t, os.RemoveAll("config.toml"))
	ctx, cancel := context.WithCancel(context.Background())
	err := watchConfig(ctx, ".")
	assert.Equal(t, "error loading config.toml: open config.toml: no such file or directory", err.Error())
	cancel()
}
