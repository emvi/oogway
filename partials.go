package oogway

import (
	"io"
	"os"
	"path/filepath"
)

const (
	partialsDir = "partials"
)

func loadPartials(dir string) error {
	d := filepath.Join(dir, partialsDir)

	if _, err := os.Stat(d); os.IsNotExist(err) || isEmptyDir(d) {
		return nil
	}

	var err error
	tpl, err = tpl.ParseGlob(filepath.Join(dir, partialsDir) + "/*.html")
	return err
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
