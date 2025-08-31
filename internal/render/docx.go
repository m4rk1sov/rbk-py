package render

import (
	"github.com/nguyenthenguyen/docx"
	"io"
)

//TODO reading docx file
func DocxReplace(tpl io.Reader, values map[string]string) ([]byte, error) {
	tmp, err := docx.ReadDocxFromMemory(tpl, 1)
}
