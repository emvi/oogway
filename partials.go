package oogway

import (
	"context"
	"html/template"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/rjeczalik/notify"
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

	change := make(chan notify.EventInfo, 1)

	go func() {
		for {
			select {
			case <-change:
				if err := loadPartials(dir, funcMap); err != nil {
					log.Printf("Error updating partials: %s", err)
				}
			case <-ctx.Done():
				notify.Stop(change)
				return
			}
		}
	}()

	if err := notify.Watch(filepath.Join(dir, partialsDir, "..."), change, notify.All); err != nil {
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
