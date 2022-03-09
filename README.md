# Oogway

Oogway is a simple web server and framework with dynamic content generation using the Go template syntax.
It's somewhere in between a static site generator and building a website manually.
Templates are automatically updated and JavaScript/TypeScript and Sass can be compiled on the fly, allowing for a quick local dev experience.
Oogway can also be used as a library in your Go application to add template functions and custom behaviour.

## Installation and Setup

Download the latest release for your platform from the release section on GitHub.
Move the binary to a directory in your $PATH (like `/usr/local/bin`).
After that, you can call Oogway from the command line using the `oogway` command.

## Configuration

Oogway is configured using a single `config.toml` file in the project directory.

```toml
[server]
host = "localhost" # leave empty for production
port = 8080
shutdown_time = 30 # time before the server is forcefully shut down (optional)
write_timeout = 5 # request write timeout
read_timeout = 5 # request read timeout

[content]
not_found = "/not-found" # specifies the redirect path when a page is not found

# optional configuration to compile sass
[sass]
dir = "assets" # asset directory path
entrypoint = "style.scss" # main sass file
out = "assets/style.css" # compiled output css file path
out_source_map = "assets/style.css.map" # css map file (optional)
watch = true # re-compile files when changed

# optional configuration to compile js/ts (see sass configuration for reference)
[js]
dir = "assets"
entrypoint = "entrypoint.js"
out = "assets/bundle.js"
out_source_map = "assets/bundle.js.map"
watch = true
```

After you have configured your project, you can start the server by running the `oogway` command inside the project directory, or by passing the directory path as the first parameter (like `oogway projects/website`).

## Structuring Your Website

There are three directories that need to be created next to the `config.toml`.

* `assets` for static files, like CSS, JavaScript, or images
* `content` for the page content and routes
* `partials` for template files that are used on multiple pages

The structure in `content` is used to create routes. Each page lives inside an `index.html`.
The home page is specified directly in the `content` directory.
Child directories can be reached by their directory name. `content/about/index.html` for example will be available on `/about`.
You can place other files next to the page to use them to build your content.
Like a markdown file which will be rendered on the page, for example.

For a demo, check out the `demo` directory on the GitHub repository.

## Template Functions

Oogway comes with a bunch of template functions that can be used to build pages.

| Function | Description | Example |
| - | - | - |
| config | Exposes the Oogway configuration. | `{{config.Server.Host}}` |
| content | Renders a template for given data. Use the route for the template name | `{{content "/about" .}}` |
| partial | Renders a partial template for given data. Use the filename without the file extension. | `{{partial "head" .}}` |
| markdown | Renders given markdown file as HTML using Go text templates. Use the full path for the template name. | `{{markdown "content/blog/article.md" .}}` |

For more functions, check out the [Sprig documentation](github.com/Masterminds/sprig).

## Using Oogway as a Library

Oogway is designed to be used as a standalone server but also as a library.
You can add your own template functions for more advanced functionality and use-cases and embed them into your application.
Simply `go get` it and call it anywhere in your application to boot up a web server.

```
import (
	"github.com/emvi/oogway"
	
	// other imports...
)

// Define a custom FuncMap to load and render blog articles from an external source.
var customFuncMap = template.FuncMap{
    "blogArticle": loadAndRenderBlogArticle,
}

func main() {
    // Start Oogway from the content/dir directory and pass your own template.FuncMap.
    // The FuncMap will be merged with the default FuncMap of Oogway.
	if err := oogway.Start("content/dir", customFuncMap); err != nil {
		log.Printf("Error starting Oogway: %s", err)
	}
}
```

## License

MIT
