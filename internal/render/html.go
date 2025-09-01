package render

import (
	"fmt"
	"github.com/cbroglie/mustache"
	"github.com/flosch/pongo2/v6"
	"github.com/m4rk1sov/rbk-py/internal/config"
	"github.com/m4rk1sov/rbk-py/internal/models"
	"path/filepath"
)

func Mustache(tpl string, data map[string]any) (string, error) {
	return mustache.Render(tpl, data)
}

func renderHTMLTemplate(cfg config.Config, req models.GenerateRequest) (string, error) {
	tplPath := filepath.Join(cfg.TemplateDir, req.Code+".html")
	htmlTpl, err := pongo2.FromFile(tplPath)
	if err != nil {
		return "", fmt.Errorf("HTML template not found: %s", req.Code)
	}

	htmlOut, err := htmlTpl.Execute(req.Data)
	if err != nil {
		return "", fmt.Errorf("template execution failed: %w", err)
	}

	return htmlOut, nil
}
