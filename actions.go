package main

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func filterOut(path string, cfg config, info os.FileInfo) bool {
	if info.IsDir() || info.Size() < cfg.size {
		return true
	}

	var isFilter bool = false
	if len(cfg.ext) > 0 {
		isFilter = true
	}
	modDate := info.ModTime()

	for _, e := range cfg.ext {
		// the ext list can never have empty strings
		// but this case is still covered for testing purposes
		if e == "" || e == filepath.Ext(path) {
			isFilter = false
		}
	}

	if !isFilter || cfg.dateAfter != "" {
		isFilter = filterDate(cfg.dateAfter, modDate)
	}

	return isFilter
}

// should return false if `filterDate` is before the `fileModTime`
func filterDate(filterDate string, fileModTime time.Time) bool {
	var filter bool = false
	if filterDate == "" {
		return filter
	}
	parsedDate, err := time.Parse(LayoutDate, filterDate)
	if err != nil {
		fmt.Println("Error formating or parsing the date can't filter", err)
		return false
	}

	if parsedDate.After(fileModTime) {
		filter = true
	}

	return filter
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
	if err := createDestDir(destDir); err != nil {
		return err
	}
	// this will get the directory where the file exists relative to the root
	relDir, err := filepath.Rel(root, filepath.Dir(path))
	if err != nil {
		return err
	}
	// this will append the .gz extension to the filename
	dest := fmt.Sprintf("%s.gz", filepath.Base(path))
	// this is creating the path to where the archive will be saved
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

func createDestDir(dirPath string) error {
	info, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(dirPath, 0755); err != nil {
				return errors.New("error creating the directory")
			}
			return nil
		}
		return err
	}

	if !info.IsDir() {
		return errors.New("please provide a valid directory name")
	}
	return nil
}
