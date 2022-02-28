package oogway

import (
	"github.com/bep/godartsass"
	"log"
	"os"
	"path/filepath"
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

// TODO call + watch
func compileSass() {
	if sass != nil {
		in := filepath.Join(cfg.Sass.Dir, cfg.Sass.Entrypoint)
		content, err := os.ReadFile(in)

		if err != nil {
			log.Printf("Error loading sass file '%s': %s", in, err)
			return
		}

		result, err := sass.Execute(godartsass.Args{
			Source:          string(content),
			OutputStyle:     godartsass.OutputStyleCompressed,
			EnableSourceMap: cfg.Sass.SourceMap,
		})

		if err != nil {
			log.Printf("Error compiling sass: %s", err)
			return
		}

		if err := os.WriteFile(cfg.Sass.Out, []byte(result.CSS), 0644); err != nil {
			log.Printf("Error writing css file '%s': %s", cfg.Sass.Out, err)
			return
		}

		if cfg.Sass.SourceMap {
			if err := os.WriteFile(cfg.Sass.OutSourceMap, []byte(result.SourceMap), 0644); err != nil {
				log.Printf("Error writing source map file '%s': %s", cfg.Sass.Out, err)
			}
		}
	}
}
