package oogway

import (
	"errors"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
)

const (
	mainScss = `body {
	max-width: 800px;
	margin: auto 0;
}
`
	mainJs = `console.log("Hello from Oogway!");
`
	indexHtml = `<!DOCTYPE html>
<html lang="en">
<head>
    <base href="/" />
    <meta charset="UTF-8" />
    <link rel="stylesheet" type="text/css" href="assets/css/main.css" />
    <script type="text/javascript" src="assets/js/main.min.js"></script>
    <title>Welcome to Oogway</title>
</head>
<body>
    <h1>Welcome to Oogway!</h1>
</body>
</html>
`
	notFoundHtml = `<!DOCTYPE html>
<html lang="en">
<head>
    <base href="/" />
    <meta charset="UTF-8" />
    <link rel="stylesheet" type="text/css" href="assets/css/main.css" />
    <script type="text/javascript" src="assets/js/main.min.js"></script>
    <title>Page not found</title>
</head>
<body>
    <h1>Page not found!</h1>
	<p>
		<a href="/">Back to home page</a>
	</p>
</body>
</html>
`
	notFoundToml = `sitemap_priority = 0.1
`
)

var (
	dirs = []string{
		"assets",
		"assets/scss",
		"assets/js",
		"content",
		"content/not-found",
		"partials",
	}
)

// Init initializes Oogway inside the specified directory.
func Init(path string) error {
	s, err := os.Stat(path)

	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0744); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	if s != nil && !s.IsDir() {
		return errors.New("target path is not a directory")
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(path, dir), 0744); err != nil {
			return err
		}
	}

	cfgFile, err := os.OpenFile(filepath.Join(path, configFile), os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer cfgFile.Close()
	cfg := Config{
		Server: ServerConfig{
			Host:            "localhost",
			Port:            8080,
			ShutdownTimeout: 30,
			WriteTimeout:    5,
			ReadTimeout:     5,
		},
		Content: ContentConfig{
			NotFound: "/not-found",
		},
		Sass: SassConfig{
			Entrypoint:   "main.scss",
			Dir:          "assets/scss",
			Watch:        true,
			Out:          "assets/css/main.css",
			OutSourceMap: "assets/css/main.css.map",
		},
		JS: JSConfig{
			Entrypoint: "main.js",
			Dir:        "assets/js",
			Watch:      true,
			Out:        "assets/js/main.min.js",
			SourceMap:  true,
		},
	}

	encoder := toml.NewEncoder(cfgFile)

	if err := encoder.Encode(&cfg); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(path, "assets", "scss", "main.scss"), []byte(mainScss), 0644); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(path, "assets", "js", "main.js"), []byte(mainJs), 0644); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(path, "content", "index.html"), []byte(indexHtml), 0644); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(path, "content", "not-found", "index.html"), []byte(notFoundHtml), 0644); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(path, "content", "not-found", "meta.toml"), []byte(notFoundToml), 0644); err != nil {
		return err
	}

	return nil
}
