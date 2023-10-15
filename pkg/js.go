package pkg

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/fsnotify/fsnotify"
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
			watcher, err := fsnotify.NewWatcher()

			if err != nil {
				return err
			}

			go func() {
				out := filepath.Join(dir, cfg.JS.Out)

				for {
					select {
					case event, ok := <-watcher.Events:
						if !ok {
							continue
						}

						if event.Op == fsnotify.Write && event.Name != out {
							ext := strings.ToLower(filepath.Ext(event.Name))

							if ext == ".js" || ext == ".ts" || ext == ".tsx" || ext == ".mts" || ext == ".cts" {
								compileJS(dir)
							}
						}
					case <-ctx.Done():
						watcher.Close()
						return
					}
				}
			}()

			if err := watcher.Add(filepath.Join(dir, cfg.JS.Dir)); err != nil {
				return err
			}
		}
	}

	return nil
}
