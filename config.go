package oogway

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/rjeczalik/notify"
)

const (
	configFile = "config.toml"
)

var (
	cfg Config
)

// Config is the Oogway application config.
type Config struct {
	Host             string
	Port             int
	ShutdownTimeout  int
	HTTPWriteTimeout int
	HTTPReadTimeout  int
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
	if cfg.Port == 0 {
		cfg.Port = 8080
	}

	if cfg.ShutdownTimeout == 0 {
		cfg.ShutdownTimeout = 30
	}

	if cfg.HTTPWriteTimeout == 0 {
		cfg.HTTPWriteTimeout = 5
	}

	if cfg.HTTPReadTimeout == 0 {
		cfg.HTTPReadTimeout = 5
	}
}

func watchConfig(ctx context.Context, dir string) error {
	if err := loadConfig(dir); err != nil {
		return err
	}

	change := make(chan notify.EventInfo, 1)

	go func() {
		for {
			select {
			case <-change:
				if err := loadConfig(dir); err != nil {
					log.Printf("Error updating config.toml: %s", err)
				}
			case <-ctx.Done():
				notify.Stop(change)
				return
			}
		}
	}()

	if err := notify.Watch(filepath.Join(dir, configFile), change, notify.Write); err != nil {
		return err
	}

	return nil
}
