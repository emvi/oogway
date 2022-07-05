package oogway

import (
	"context"
	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/rjeczalik/notify"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func compileJS(dir string) {
	in := filepath.Join(dir, cfg.JS.Dir, cfg.JS.Entrypoint)
	log.Printf("Compiling js file: %s", in)
	sourceMap := esbuild.SourceMapNone

	if cfg.JS.SourceMap {
		sourceMap = esbuild.SourceMapExternal
	}

	if err := os.MkdirAll(filepath.Join(dir, filepath.Dir(cfg.JS.Out)), 0744); err != nil {
		log.Printf("Error creating js output directory: %s", err)
		return
	}

	result := esbuild.Build(esbuild.BuildOptions{
		EntryPoints:       []string{in},
		Outfile:           filepath.Join(dir, cfg.JS.Out),
		Sourcemap:         sourceMap,
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Write:             true,
	})

	if len(result.Errors) > 0 {
		log.Println("Error compiling js:")

		for _, err := range result.Errors {
			log.Println(err.Text)
		}
	}
}

func watchJS(ctx context.Context, dir string) error {
	if cfg.JS.Entrypoint != "" {
		compileJS(dir)

		if cfg.JS.Watch {
			bundle, err := filepath.Abs(filepath.Join(dir, cfg.JS.Out))

			if err != nil {
				return err
			}

			change := make(chan notify.EventInfo, 1)

			go func() {
				for {
					select {
					case event := <-change:
						if event.Path() != bundle {
							ext := strings.ToLower(filepath.Ext(event.Path()))

							if ext == ".js" || ext == ".ts" || ext == ".tsx" || ext == ".mts" || ext == ".cts" {
								compileJS(dir)
							}
						}
					case <-ctx.Done():
						notify.Stop(change)
						return
					}
				}
			}()

			if err := notify.Watch(filepath.Join(dir, cfg.JS.Dir, "..."), change, notify.Write); err != nil {
				return err
			}
		}
	}

	return nil
}
