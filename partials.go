package oogway

import (
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

const (
	partialsDir = "partials"
	tplFileExt  = ".html"
)

func loadPartials(dir string) error {
	d := filepath.Join(dir, partialsDir)

	if _, err := os.Stat(d); os.IsNotExist(err) || isEmptyDir(d) {
		return nil
	}

	return filepath.WalkDir(d, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && filepath.Ext(path) == tplFileExt {
			if _, err := tpl.load(path); err != nil {
				log.Printf("Error loading template %s: %s", path, err)
			}
		}

		return nil
	})
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
