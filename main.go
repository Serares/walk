package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type config struct {
	// extension to filter
	ext string
	// min file size
	size int64
	// list files
	list bool
	// delete files
	del bool
	// log destination writter
	wLog io.Writer
	// archive dir
	archive string
}

func main() {
	archive := flag.String("archive", "", "Archive directory")
	root := flag.String("root", "", "The directory where to start crawling")
	list := flag.Bool("list", false, "List files only")
	ext := flag.String("ext", "", "File extension to search for")
	size := flag.Int64("size", 0, "Minimum file size")
	del := flag.Bool("del", false, "Delete files")
	logFile := flag.String("log", "", "Log deletes to the specified file")
	flag.Parse()

	var (
		f   = os.Stdout
		err error
	)

	if *logFile != "" {
		f, err = os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()
	}

	c := config{
		archive: *archive,
		ext:     *ext,
		size:    *size,
		list:    *list,
		del:     *del,
		wLog:    f,
	}

	if err := run(*root, f, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(root string, w io.Writer, cfg config) error {
	delLogger := log.New(cfg.wLog, "DELETED FILE:", log.LstdFlags)
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filterOut(path, cfg.ext, cfg.size, info) {
			return nil
		}

		if cfg.list {
			return listFile(path, w)
		}

		if cfg.archive != "" {
			if err := archiveFile(cfg.archive, root, path); err != nil {
				return err
			}
		}

		if cfg.del {
			return delFile(path, delLogger)
		}

		return listFile(path, w)
	})
}
