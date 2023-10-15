package oogway

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/fsnotify/fsnotify"
)

const (
	configFile = "config.toml"
)

var (
	cfg Config
)

// Config is the Oogway application config.
type Config struct {
	Server  ServerConfig  `toml:"server"`
	Content ContentConfig `toml:"content"`
	Sass    SassConfig    `toml:"sass"`
	JS      JSConfig      `toml:"js"`
	Pirsch  PirschConfig  `toml:"pirsch"`
}

// ServerConfig is the HTTP server configuration.
type ServerConfig struct {
	Host            string `toml:"host"`
	Port            int    `toml:"port"`
	ShutdownTimeout int    `toml:"shutdown_time"`
	WriteTimeout    int    `toml:"write_timeout"`
	ReadTimeout     int    `toml:"read_timeout"`
	TLSCertFile     string `toml:"tls_cert_file"`
	TLSKeyFile      string `toml:"tls_key_file"`
}

// ContentConfig is the content configuration.
type ContentConfig struct {
	NotFound string `toml:"not_found"`
}

// SassConfig is the sass compiler configuration.
type SassConfig struct {
	Entrypoint   string `toml:"entrypoint"`
	Dir          string `toml:"dir"`
	Watch        bool   `toml:"watch"`
	Out          string `toml:"out"`
	OutSourceMap string `toml:"out_source_map"`
}

// JSConfig is the JavaScript compiler configuration.
type JSConfig struct {
	Entrypoint string `toml:"entrypoint"`
	Dir        string `toml:"dir"`
	Watch      bool   `toml:"watch"`
	Out        string `toml:"out"`
	SourceMap  bool   `toml:"source_map"`
}

// PirschConfig is the configuration for pirsch.io.
type PirschConfig struct {
	ClientID     string   `toml:"client_id"`
	ClientSecret string   `toml:"client_secret"`
	Subnets      []string `toml:"subnets"`
	Header       []string `toml:"header"`
}

func loadConfig(dir string) error {
	content, err := os.ReadFile(filepath.Join(dir, configFile))

	if err != nil {
		return fmt.Errorf("error loading config.toml: %s", err)
	}

	if _, err := toml.Decode(string(content), &cfg); err != nil {
		return fmt.Errorf("error loading config.toml: %s", err)
	}

	setConfigDefaults()
	return nil
}

func setConfigDefaults() {
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}

	if cfg.Server.ShutdownTimeout == 0 {
		cfg.Server.ShutdownTimeout = 30
	}

	if cfg.Server.WriteTimeout == 0 {
		cfg.Server.WriteTimeout = 5
	}

	if cfg.Server.ReadTimeout == 0 {
		cfg.Server.ReadTimeout = 5
	}
}

func watchConfig(ctx context.Context, dir string) error {
	if err := loadConfig(dir); err != nil {
		return err
	}

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					continue
				}

				if event.Op == fsnotify.Write {
					if err := loadConfig(dir); err != nil {
						log.Printf("Error updating config.toml: %s", err)
					}
				}
			case <-ctx.Done():
				watcher.Close()
				return
			}
		}
	}()

	if err := watcher.Add(filepath.Join(dir, configFile)); err != nil {
		return err
	}

	return nil
}
