package oogway

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	sampleConfig = `hostname = "localhost"
port = 8080`
)

func TestLoadConfig(t *testing.T) {
	assert.NoError(t, os.RemoveAll("config.toml"))
	assert.NoError(t, os.WriteFile("config.toml", []byte(sampleConfig), 0644))
	assert.NoError(t, loadConfig("."))
	assert.Equal(t, "localhost", cfg.Hostname)
	assert.Equal(t, 8080, cfg.Port)
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
	assert.Equal(t, 8080, cfg.Port)
	assert.NoError(t, os.WriteFile("config.toml", []byte(strings.Replace(sampleConfig, "8080", "8888", 1)), 0644))
	time.Sleep(time.Millisecond * 10)
	cancel()
	assert.Equal(t, 8888, cfg.Port)
}

func TestWatchConfigNotExists(t *testing.T) {
	assert.NoError(t, os.RemoveAll("config.toml"))
	ctx, cancel := context.WithCancel(context.Background())
	err := watchConfig(ctx, ".")
	assert.Equal(t, "error loading config.toml: open config.toml: no such file or directory", err.Error())
	cancel()
}
