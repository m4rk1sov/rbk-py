package httpserver

import (
	"fmt"
	"net/http"
	"path/filepath"
	
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/m4rk1sov/rbk-py/internal/config"
)

func handleGenerateHTML(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.GenerateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			web.BadRequest(c, "invalid JSON: "+err.Error())
			return
		}
		if req.Code == "" {
			web.BadRequest(c, "code is required")
			return
		}
		// read .html
		tplPath := filepath.Join(cfg.TemplateDir, req.Code+".html")
		htmlTpl, err := templates.ReadFile(tplPath)
		if err != nil {
			web.Error(c, http.StatusNotFound, fmt.Errorf("template not found: %s", req.Code))
			return
		}
		// logic for rendering html??
		htmlOut, err := render.HTMLRender(htmlTpl, req.Data)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, err)
			return
		}
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Header("Content-Disposition", fmt.Sprintf(`inline; filename="%s.html"`, req.Code))
		c.String(http.StatusOK, htmlOut)
	}
}
