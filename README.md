# Oogway

Oogway is a simple web server and framework with dynamic content generation using the Go template syntax.
It's somewhere between a static site generator and manually building a website.
Templates are updated automatically and JavaScript/TypeScript and Sass can be compiled on the fly, allowing for a fast local development experience.
Oogway can also be used as a library in your Go application to add template functionality and custom behavior.

## Installation and Setup

Download the latest release for your platform from the releases section on GitHub.
Move the binary to a directory in your $PATH (like `/usr/local/bin`).
For Sass, you need to install the `sass` command globally (`sudo npm i -g sass`).
After that you can run Oogway from the command line with the `oogway` command.

* `oogway run <path>` will run Oogway in the given directory.
* `oogway init <path>` will initialize a new project in the given directory.

Or through Docker:

```yaml
version: "3"

services:
  oogway:
    image: emvicom/oogway
    container_name: oogway
    restart: unless-stopped
    ports:
      - "8080:8080"
    volumes:
      - ./demo:/app/data
```

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

# optional configuration for pirsch.io
[pirsch]
client_id = "..." # optional when using an access key (recommended) instead of oAuth
client_secret = "..." # required
```

After you have configured your project, you can start the server by running the `oogway` command inside the project directory, or by passing the directory path as the first parameter (like `oogway projects/website`).

## Structuring Your Website

There are three directories that need to be created in addition to `config.toml

* `assets` for static files like CSS, JavaScript or images
* `content` for page content and routes
* `partials` for template files used on multiple pages

The structure in `content` is used to create routes. Each page lives within an `index.html`.
The start page is specified directly in the `content` directory.
Subdirectories can be accessed by their directory names. For example, `content/about/index.html` will be accessible from `/about`.
You can place other files next to the page to use in building your content.
For example, a markdown file that is rendered on the page.

A `meta.toml` file can be created next to each `index.html` for additional configuration.

```toml
# sets the priority in the sitemap.xml. Default is 1
sitemap_priority = 0.95
```

You can find a demo in the `demo` directory of the GitHub repository.

## Template Functions

Oogway comes with a number of template functions that can be used to create pages.

| Function      | Description                                                                                                        | Example                                                     |
|---------------|--------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------|
| config        | Exposes the Oogway configuration.                                                                                  | `{{config.Server.Host}}`                                    |
| content       | Renders a template for given data. Use the route for the template name                                             | `{{content "/about" .}}`                                    |
| partial       | Renders a partial template for given data. Use the filename without the file extension.                            | `{{partial "head" .}}`                                      |
| markdown      | Renders given markdown file as HTML using Go text templates. Use the full path for the template name.              | `{{markdown "content/blog/article.md" .}}`                  |
| markdownBlock | Renders a block from given markdown file as HTML using Go text templates. Use the full path for the template name. | `{{markdownBlock "content/blog/article.md" "blockName" .}}` |
| int           | Converts given string to an integer.                                                                               | `{{int "123"}}`                                             |
| uint64        | Converts given int to an uint64.                                                                                   | `{{uint64 123}}`                                            |

For more features, see the [Sprig documentation](github.com/Masterminds/sprig).

## Using Oogway as a Library

Oogway is designed to be used as a standalone server, but also as a library.
You can add your own template functions for more advanced functionality and use cases and embed them in your application.
Just `go get` it and call it anywhere in your application to start a web server.

```
import (
	oogway "github.com/emvi/oogway/pkg"
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
