package render

import "github.com/cbroglie/mustache"

func Mustache(tpl string, data map[string]any) (string, error) {
	return mustache.Render(tpl, data)
}
