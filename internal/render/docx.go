package render

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lukasjarosch/go-docx"
	"github.com/m4rk1sov/rbk-py/internal/config"
	"github.com/m4rk1sov/rbk-py/internal/models"
	"github.com/m4rk1sov/rbk-py/templates"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

func GenerateDocxFromTemplate(cfg config.Config, req models.GenerateRequest) ([]byte, error) {
	tplPath := filepath.Join(cfg.TemplateDir, req.Code+".docx")
	buf, err := templates.ReadBinaryFile(tplPath)
	if err != nil {
		return nil, fmt.Errorf("template not found: %s", req.Code)
	}

	m, flatErr := models.FlattenToStringMap(req.Data, "")
	if flatErr != nil {
		return nil, flatErr
	}

	out, err := docxReplace(bytes.NewReader(buf), m)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func GenerateDocxFromHTML(c *gin.Context, cfg config.Config, req models.GenerateRequest) ([]byte, error) {
	// 1. Render HTML template
	htmlContent, err := renderHTMLTemplate(cfg, req)
	if err != nil {
		return nil, fmt.Errorf("failed to render HTML template: %w", err)
	}

	// 2. Convert HTML to DOCX using external service
	docxBytes, err := htmlToDocx(c, cfg.PDFConverterURL, req.Code+".html", strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("HTML to DOCX conversion failed: %w", err)
	}

	return docxBytes, nil
}

// htmlToDocx converts HTML content to DOCX using external conversion service
func htmlToDocx(ctx context.Context, converterURL, filename string, htmlReader io.Reader) ([]byte, error) {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("files", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = io.Copy(part, htmlReader)
	if err != nil {
		return nil, fmt.Errorf("failed to copy HTML content: %w", err)
	}

	err = writer.WriteField("pdfFormat", "PDF/A-1a")
	if err != nil {
		return nil, err
	}
	err = writer.WriteField("pdfUniversalAccess", "true")
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		converterURL+"/forms/libreoffice/convert",
		&requestBody,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("conversion request failed: %w", err)
	}
	defer func(Body io.ReadCloser) {
		if closeErr := Body.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("conversion failed with status %d: %s", resp.StatusCode, string(body))
	}

	docxBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read converted DOCX: %w", err)
	}

	return docxBytes, nil
}

// Alternative implementation using Pandoc (if you prefer local conversion)
func HTMLToDocxWithPandoc(htmlContent string) ([]byte, error) {
	// This would require Pandoc to be installed on the system
	// Implementation would use exec.Command to call pandoc
	return nil, fmt.Errorf("pandoc conversion not implemented yet")
}

// Alternative implementation using LibreOffice headless
func HTMLToDocxWithLibreOffice(htmlFilePath string) ([]byte, error) {
	// This would require LibreOffice to be installed on the system
	// Implementation would use exec.Command to call libreoffice --headless
	return nil, fmt.Errorf("libreoffice conversion not implemented yet")
}

//// docxReplace replaces placeholders in a DOCX template using go-docx.
//// - `values` can contain both simple strings and slices of maps (for loops).
//func docxReplace(tpl io.Reader, values map[string]string) ([]byte, error) {
//	// Load the template from memory (requires a []byte)
//	buf := new(bytes.Buffer)
//	if _, err := io.Copy(buf, tpl); err != nil {
//		return nil, fmt.Errorf("failed to read template: %w", err)
//	}
//
//	// Open DOCX template
//	doc, err := docx.OpenBytes(buf.Bytes())
//	if err != nil {
//		return nil, fmt.Errorf("failed to open DOCX template: %w", err)
//	}
//
//	// Replace simple values
//	if err := doc.ReplaceAll(values); err != nil {
//		return nil, fmt.Errorf("failed to replace placeholders: %w", err)
//	}
//
//	// Handle loops: detect keys with []map[string]interface{}
//	for key, v := range values {
//		if list, ok := v.([]map[string]interface{}); ok {
//			if err := doc.ReplaceLoop(key, list); err != nil {
//				return nil, fmt.Errorf("failed to replace loop %q: %w", key, err)
//			}
//		}
//	}
//
//	// Write back into memory
//	var out bytes.Buffer
//	if err := doc.Write(&out); err != nil {
//		return nil, fmt.Errorf("failed to write final DOCX: %w", err)
//	}
//
//	return out.Bytes(), nil
//}

func docxReplace(tpl io.Reader, values map[string]string) ([]byte, error) {
	src, err := io.ReadAll(tpl)
	if err != nil {
		return nil, fmt.Errorf("read docx template: %w", err)
	}

	doc, err := docx.OpenBytes(src)
	if err != nil {
		return nil, fmt.Errorf("open docx: %w", err)
	}
	defer doc.Close()

	// Expand the map so both {{key}} and {key} are recognized
	expanded := make(map[string]any, len(values)*2)
	for k, v := range values {
		expanded[k] = v
		expanded["{{"+k+"}}"] = v
		expanded["{{ "+k+" }}"] = v // with spaces
	}

	if err := doc.ReplaceAll(docx.PlaceholderMap(expanded)); err != nil {
		return nil, fmt.Errorf("replace placeholders: %w", err)
	}

	var out bytes.Buffer
	if err := doc.Write(&out); err != nil {
		return nil, fmt.Errorf("write result docx: %w", err)
	}
	return out.Bytes(), nil
}
