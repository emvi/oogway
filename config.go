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
	cfg config
)

type config struct {
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

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					watcher.Close()
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					if err := loadConfig(dir); err != nil {
						log.Printf("error updating config.toml: %s", err)
					}
				}
			case err, ok := <-watcher.Errors:
				watcher.Close()

				if !ok {
					return
				}

				panic(err)
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
