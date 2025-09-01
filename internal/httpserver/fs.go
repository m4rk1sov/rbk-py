package httpserver

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func List(dir string) ([]string, error) {
	entries := map[string]struct{}{}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		name := filepath.Base(path)
		ext := filepath.Ext(name)
		base := strings.TrimSuffix(name, ext)
		entries[base] = struct{}{}
		return nil
	})
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, err
	}
	out := make([]string, 0, len(entries))
	for k := range entries {
		out = append(out, k)
	}
	return out, nil
}

func ReadFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	return string(b), err
}

func ReadBinaryFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
