package oogway

import (
	"context"
	"html/template"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

const (
	partialsDir = "partials"
	tplFileExt  = ".html"
)

var (
	partials = newTplCache()
)

func loadPartials(dir string, funcMap template.FuncMap) error {
	partials.clear()
	d := filepath.Join(dir, partialsDir)

	if _, err := os.Stat(d); os.IsNotExist(err) || isEmptyDir(d) {
		return nil
	}

	return filepath.WalkDir(d, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && filepath.Ext(path) == tplFileExt {
			if _, err := partials.load(path, funcMap); err != nil {
				log.Printf("Error loading template %s: %s", path, err)
			}
		}

		return nil
	})
}

func watchPartials(ctx context.Context, dir string, funcMap template.FuncMap) error {
	if err := loadPartials(dir, funcMap); err != nil {
		return err
	}

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case _, ok := <-watcher.Events:
				if !ok {
					continue
				}

				if err := loadPartials(dir, funcMap); err != nil {
					log.Printf("Error updating partials: %s", err)
				}
			case <-ctx.Done():
				watcher.Close()
				return
			}
		}
	}()

	if err := watcher.Add(filepath.Join(dir, partialsDir)); err != nil {
		return err
	}

	return nil
}

func isEmptyDir(path string) bool {
	f, err := os.Open(path)

	if err != nil {
		return true
	}

	defer f.Close()
	_, err = f.Readdirnames(1)
	return err == io.EOF
}
