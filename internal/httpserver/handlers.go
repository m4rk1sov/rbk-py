package httpserver

import (
	"fmt"
	"github.com/flosch/pongo2/v6"
	"github.com/m4rk1sov/rbk-py/internal/models"
	"github.com/m4rk1sov/rbk-py/internal/pdf"
	"github.com/m4rk1sov/rbk-py/internal/render"
	"github.com/m4rk1sov/rbk-py/internal/web"
	"github.com/m4rk1sov/rbk-py/templates"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
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
			if req.Code == "" {
				web.BadRequest(c, "code is required")
				return
			}
			// read .html
			tplPath := filepath.Join(cfg.TemplateDir, req.Code+".html")
			htmlTpl, err := pongo2.FromFile(tplPath)
			if err != nil {
				web.Error(c, http.StatusNotFound, fmt.Errorf("template not found: %s", req.Code))
				return
			}
			htmlOut, err := htmlTpl.Execute(req.Data)
			if err != nil {
				web.Error(c, http.StatusInternalServerError, err)
				return
			}
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.Header("Content-Disposition", fmt.Sprintf(`inline; filename="%s.html"`, req.Code))
			c.String(http.StatusOK, htmlOut)
		}
	}
}

func handleGeneratePDF(cfg config.Config) gin.HandlerFunc {
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
		htmlTpl, err := pongo2.FromFile(tplPath)
		if err != nil {
			web.Error(c, http.StatusNotFound, fmt.Errorf("template not found: %s", req.Code))
			return
		}
		htmlOut, err := htmlTpl.Execute(req.Data)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, err)
			return
		}

		pdfBytes, err := pdf.HTMLToPDF(c, cfg.PDFConverterURL, req.Code+".html", strings.NewReader(htmlOut))
		if err != nil {
			web.Error(c, http.StatusBadGateway, fmt.Errorf("pdf conversion failed: %w", err))
			return
		}
		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.pdf"`, req.Code))
		c.Data(http.StatusOK, "application/pdf", pdfBytes)
	}
}

func handleGenerateDOCX(cfg config.Config) gin.HandlerFunc {
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

		htmlTplPath := filepath.Join(cfg.TemplateDir, req.Code+".html")
		if _, err := templates.ReadBinaryFile(htmlTplPath); err == nil {
			// HTML template exists, use hybrid approach
			docxBytes, err := render.GenerateDocxFromHTML(c, cfg, req)
			if err != nil {
				// Fall back to traditional DOCX template if HTML conversion fails
				fallbackDocxBytes, fallbackErr := render.GenerateDocxFromTemplate(cfg, req)
				if fallbackErr != nil {
					web.Error(c, http.StatusInternalServerError, fmt.Errorf("both HTML conversion and DOCX template failed: %w, fallback: %w", err, fallbackErr))
					return
				}
				docxBytes = fallbackDocxBytes
			}

			fn := req.Code + ".docx"
			ct := "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
			c.Header("Content-Type", ct)
			c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fn))
			c.Data(http.StatusOK, ct, docxBytes)
			return
		}

		// Fall back to traditional DOCX template approach
		docxBytes, err := render.GenerateDocxFromTemplate(cfg, req)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, err)
			return
		}

		//tplPath := filepath.Join(cfg.TemplateDir, req.Code+".docx")
		//buf, err := templates.ReadBinaryFile(tplPath)
		//if err != nil {
		//	web.Error(c, http.StatusNotFound, fmt.Errorf("template not found: %s", req.Code))
		//	return
		//}

		//m, flatErr := models.FlattenToStringMap(req.Data, "")
		//if flatErr != nil {
		//	web.Error(c, http.StatusBadRequest, flatErr)
		//	return
		//}
		//out, err := render.DocxReplace(bytes.NewReader(buf), m)
		//if err != nil {
		//	web.Error(c, http.StatusInternalServerError, err)
		//	return
		//}

		fn := req.Code + ".docx"
		ct := mime.TypeByExtension(".docx")
		if ct == "" {
			ct = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
		}
		c.Header("Content-Type", ct)
		c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fn))
		c.Data(http.StatusOK, ct, docxBytes)

		//_, err = io.Copy(c.Writer, bytes.NewReader(out))
		//if err != nil {
		//	return
		//}
	}
}

func handleGenerateXLSX(cfg config.Config) gin.HandlerFunc {
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
		tplPath := filepath.Join(cfg.TemplateDir, req.Code+".xlsx")
		tplBytes, err := templates.ReadBinaryFile(tplPath)
		if err != nil {
			web.Error(c, http.StatusNotFound, fmt.Errorf("template not found: %s", req.Code))
			return
		}
		outBytes, err := render.ExcelizeFill(tplBytes, req.Data)
		if err != nil {
			web.Error(c, http.StatusInternalServerError, err)
			return
		}

		fn := req.Code + ".xlsx"
		ct := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
		c.Header("Content-Type", ct)
		c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fn))
		c.Data(http.StatusOK, ct, outBytes)
	}
}
