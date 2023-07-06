package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
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

func delFile(path string, delLog *log.Logger) error {
	if err := os.Remove(path); err != nil {
		return err
	}
	delLog.Println(path)
	return nil
}

func archiveFile(destDir, root, path string) error {
	info, err := os.Stat(destDir)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", destDir)
	}

	// this will get the directory where the file exists relative to the root
	relDir, err := filepath.Rel(root, filepath.Dir(path))
	if err != nil {
		return err
	}
	// this will append the .gz extension to the filename
	dest := fmt.Sprintf("%s.gz", filepath.Base(path))
	// this is creating the path to where the archive willbe saved
	targetPath := filepath.Join(destDir, relDir, dest)
	// this will create all the directories at once; if they exist it will
	// do nothing
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	out, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer out.Close()
	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()
	zw := gzip.NewWriter(out)
	zw.Name = filepath.Base(path)

	if _, err = io.Copy(zw, in); err != nil {
		return err
	}
	if err := zw.Close(); err != nil {
		return err
	}

	return out.Close()
}
