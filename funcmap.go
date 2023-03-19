package oogway

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"html/template"
	"log"
	"os"
	"path/filepath"
	tt "text/template"
)

var (
	defaultFuncMap = template.FuncMap{
		"config":        func() Config { return cfg },
		"content":       renderContent,
		"partial":       renderPartial,
		"markdown":      renderMarkdown,
		"markdownBlock": renderMarkdownBlock,
	}
)

func mergeFuncMaps(maps ...template.FuncMap) template.FuncMap {
	out := make(map[string]any)

	for k, v := range sprig.FuncMap() {
		out[k] = v
	}

	for k, v := range defaultFuncMap {
		out[k] = v
	}

	for _, m := range maps {
		if m != nil {
			for k, v := range m {
				out[k] = v
			}
		}
	}

	return out
}

func renderContent(tpl string, data any) template.HTML {
	c := content.get(tpl)

	if c == nil {
		return template.HTML(fmt.Sprintf("Template '%s' not found", tpl))
	}

	var buffer bytes.Buffer

	if err := c.Execute(&buffer, data); err != nil {
		return template.HTML(fmt.Sprintf("Error rendering template '%s': %s", tpl, err))
	}

	return template.HTML(buffer.String())
}

func renderPartial(tpl string, data any) template.HTML {
	partial := partials.get(tpl)

	if partial == nil {
		return template.HTML(fmt.Sprintf("Partial '%s' not found", tpl))
	}

	var buffer bytes.Buffer

	if err := partial.Execute(&buffer, data); err != nil {
		return template.HTML(fmt.Sprintf("Error rendering partial '%s': %s", tpl, err))
	}

	return template.HTML(buffer.String())
}

func renderMarkdown(file string, data any) template.HTML {
	content, err := os.ReadFile(filepath.Join(baseDir, file))

	if err != nil {
		log.Printf("Error loading markdown file '%s': %s", file, err)
		return ""
	}

	return renderMarkdownContent(file, string(content), "", data)
}

func renderMarkdownBlock(file, block string, data any) template.HTML {
	content, err := os.ReadFile(filepath.Join(baseDir, file))

	if err != nil {
		log.Printf("Error loading markdown file '%s': %s", file, err)
		return ""
	}

	return renderMarkdownContent(file, string(content), block, data)
}

func renderMarkdownContent(file, content, block string, data any) template.HTML {
	tpl, err := tt.New("").Funcs(tplFuncMap).Parse(content)

	if err != nil {
		log.Printf("Error parsing markdown file '%s': %s", file, err)
		return ""
	}

	var buffer, out bytes.Buffer
	converter := goldmark.New(
		goldmark.WithExtensions(
			extension.NewFootnote(),
			extension.NewTable(),
			extension.Strikethrough,
			extension.TaskList,
		),
	)

	if block != "" {
		if _, err := tpl.Parse(fmt.Sprintf(`{{template "%s" .}}`, block)); err != nil {
			log.Printf("Error parsing markdown block '%s' in file '%s': %s", block, file, err)
			return ""
		}

		if err := tpl.Execute(&buffer, data); err != nil {
			log.Printf("Error rendering markdown file '%s': %s", file, err)
			return ""
		}

		if err := converter.Convert(buffer.Bytes(), &out); err != nil {
			panic(err)
		}

		return template.HTML(out.String())
	}

	if err := tpl.Execute(&buffer, data); err != nil {
		log.Printf("Error rendering markdown file '%s': %s", file, err)
		return ""
	}

	if err := converter.Convert(buffer.Bytes(), &out); err != nil {
		panic(err)
	}

	return template.HTML(out.String())
}
