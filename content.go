package oogway

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

const (
	contentDir      = "content"
	contentPageFile = "index.html"
)

var (
	tpl = newTplCache()
)

// TODO add/return routes
func loadContent(dir string) error {
	d := filepath.Join(dir, contentDir)

	if _, err := os.Stat(d); os.IsNotExist(err) || isEmptyDir(d) {
		return nil
	}

	return filepath.WalkDir(d, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && d.Name() == contentPageFile {
			if err := tpl.load(path); err != nil {
				log.Printf("Error loading template %s: %s", path, err)
			}
		}

		return nil
	})
}
