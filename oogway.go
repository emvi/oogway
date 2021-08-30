package oogway

import (
	"context"
	"html/template"
)

// Start starts the Oogway server for given directory.
// The second argument is an optional template.FuncMap that will be merged into Oogway's funcmap.
func Start(dir string, funcmap template.FuncMap) error {
	ctx, cancel := context.WithCancel(context.Background())

	if err := watchConfig(ctx, dir); err != nil {
		cancel()
		return err
	}

	cancel() // TODO
	return nil
}
