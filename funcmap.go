package oogway

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"html/template"
)

var (
	defaultFuncMap = template.FuncMap{
		"config":  func() Config { return cfg },
		"content": renderContent,
		"partial": renderPartial,
	}
)

func mergeFuncMaps(maps ...template.FuncMap) template.FuncMap {
	out := make(map[string]interface{})

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

func renderContent(tpl string, data interface{}) template.HTML {
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

func renderPartial(tpl string, data interface{}) template.HTML {
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
