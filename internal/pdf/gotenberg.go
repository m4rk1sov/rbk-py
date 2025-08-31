package pdf

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"
)

func HTMLToPDF(c *gin.Context, baseURL, fileName string, html io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	
	//TODO fix the naming of html
	fw, err := mw.CreateFormFile("files", "index.html")
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(fw, html); err != nil {
		return nil, err
	}
	
	err = mw.WriteField("paperWidth", "8.27")
	err = mw.WriteField("paperHeight", "11.69")
	err = mw.WriteField("marginTop", "0.5")
	err = mw.WriteField("paperBottom", "0.5")
	err = mw.WriteField("marginLeft", "0.5")
	err = mw.WriteField("marginRight", "0.5")
	err = mw.Close()
	if err != nil {
		return nil, err
	}
	
	url := fmt.Sprintf("%s/forms/chromium/convert/html", baseURL)
	req, err := http.NewRequestWithContext(c, http.MethodPost, url, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", mw.FormDataContentType())
	
	client := &http.Client{Timeout: 60 * time.Second}
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(resp io.ReadCloser) {
		if closeErr := resp.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}(resp.Body)
	
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gotenberg status %d: %s", resp.Status, string(b))
	}
	return io.ReadAll(resp.Body)
}

func ext(base, e string) string {
	return base + filepath.Ext(e)
}
