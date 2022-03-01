package oogway

import (
	"context"
	"github.com/bep/godartsass"
	"github.com/rjeczalik/notify"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	sass *godartsass.Transpiler
)

func init() {
	var err error
	sass, err = godartsass.Start(godartsass.Options{})

	if err != nil {
		log.Printf("Error setting up sass compiler: %s. Oogway will still work, but sass compilation won't be available.", err)
	}
}

func compileSass(dir string) {
	if sass != nil {
		in := filepath.Join(dir, cfg.Sass.Dir, cfg.Sass.Entrypoint)
		log.Printf("Compiling sass file: %s", in)
		content, err := os.ReadFile(in)

		if err != nil {
			log.Printf("Error loading sass file '%s': %s", in, err)
			return
		}

		result, err := sass.Execute(godartsass.Args{
			Source:          string(content),
			OutputStyle:     godartsass.OutputStyleCompressed,
			EnableSourceMap: cfg.Sass.OutSourceMap != "",
		})

		if err != nil {
			log.Printf("Error compiling sass: %s", err)
			return
		}

		out := filepath.Join(dir, cfg.Sass.Out)

		if err := os.WriteFile(out, []byte(result.CSS), 0644); err != nil {
			log.Printf("Error writing css file '%s': %s", out, err)
			return
		}

		if cfg.Sass.OutSourceMap != "" {
			out = filepath.Join(dir, cfg.Sass.OutSourceMap)

			if err := os.WriteFile(out, []byte(result.SourceMap), 0644); err != nil {
				log.Printf("Error writing source map file '%s': %s", out, err)
			}
		}
	}
}

func watchSass(ctx context.Context, dir string) error {
	if cfg.Sass.Entrypoint != "" {
		compileSass(dir)

		if cfg.Sass.Watch {
			change := make(chan notify.EventInfo, 1)

			go func() {
				for {
					select {
					case event := <-change:
						ext := strings.ToLower(filepath.Ext(event.Path()))

						if ext == ".scss" || ext == ".sass" {
							compileSass(dir)
						}
					case <-ctx.Done():
						notify.Stop(change)
						return
					}
				}
			}()

			if err := notify.Watch(filepath.Join(dir, cfg.Sass.Dir), change, notify.Write); err != nil {
				return err
			}
		}
	}

	return nil
}
