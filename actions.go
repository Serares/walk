package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func filterOut(path string, ext string, minSzie int64, info os.FileInfo) bool {
	if info.IsDir() || info.Size() < minSzie {
		return true
	}

	if ext != "" && filepath.Ext(path) != ext {
		return true
	}

	return false
}

func listFile(path string, w io.Writer) error {
	_, err := fmt.Fprintln(w, path)
	return err
}
