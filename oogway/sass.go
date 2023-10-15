package oogway

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func compileSass(dir string) {
	if err := os.MkdirAll(filepath.Join(dir, filepath.Dir(cfg.Sass.Out)), 0744); err != nil {
		log.Printf("Error creating css output directory: %s", err)
		return
	}

	in := filepath.Join(dir, cfg.Sass.Dir, cfg.Sass.Entrypoint)
	out := filepath.Join(dir, cfg.Sass.Out)
	log.Printf("Compiling sass file '%s' to '%s'", in, out)
	dirs, err := getDirs(filepath.Join(dir, cfg.Sass.Dir))

	if err != nil {
		log.Printf("Error reading sass directory: %s", err)
		return
	}

	args := make([]string, 0)

	for _, d := range dirs {
		args = append(args, fmt.Sprintf("--load-path=%s", d))
	}

	if cfg.Sass.OutSourceMap == "" {
		args = append(args, "--no-source-map")
	} else {
		args = append(args, "--source-map")
	}

	args = append(args, "--style=compressed")
	args = append(args, in)
	args = append(args, out)
	cmd := exec.Command("sass", args...)

	if err := cmd.Run(); err != nil {
		log.Printf("Error compiling sass: %s", err)
		return
	}
}

func watchSass(ctx context.Context, dir string) error {
	if cfg.Sass.Entrypoint != "" {
		compileSass(dir)

		if cfg.Sass.Watch {
			watcher, err := fsnotify.NewWatcher()

			if err != nil {
				return err
			}

			go func() {
				out := filepath.Join(dir, cfg.Sass.Out)

				for {
					select {
					case event, ok := <-watcher.Events:
						if !ok {
							continue
						}

						if event.Op == fsnotify.Write && event.Name != out {
							ext := strings.ToLower(filepath.Ext(event.Name))

							if ext == ".scss" || ext == ".sass" {
								compileSass(dir)
							}
						}
					case <-ctx.Done():
						watcher.Close()
						return
					}
				}
			}()

			if err := watcher.Add(filepath.Join(dir, cfg.Sass.Dir)); err != nil {
				return err
			}
		}
	}

	return nil
}

func getDirs(root string) ([]string, error) {
	dirs := make([]string, 0)

	if err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			dirs = append(dirs, path)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return dirs, nil
}
