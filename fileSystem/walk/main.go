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
	// extension to filter out
	ext string

	// min file size
	size int64

	// list files
	list bool

	// delete matching file
	del bool

	wLog io.Writer
}



func main() {
	root := flag.String("root", ".", "root directory")
	list := flag.Bool("list", false, "list files")
	ext := flag.String("ext", "", "file extension")
	logFile := flag.String("log", "", "log deletes to this file")
	size := flag.Int64("size", 0, "minimum file size")
	del := flag.Bool("del", false, "delete matching file")

	flag.Parse()

var (
	f = os.Stdout
	err error
)

if *logFile != ""{
	f, err  = os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()
}

	c := config{
		ext:  *ext,
		size: *size,
		list: *list,
		del:  *del,
		wLog: f,
	}

	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(root string, out io.Writer, c config) error {
	delLogger:= log.New(c.wLog, "DELETED FILE", log.LstdFlags)
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filterOut(path, c.ext, c.size, info) {
			return nil
		}

		// list explicity set - dont do anything else
		if c.list {
			return listFiles(path, out)
		}

		if c.del {
			return delFile(path, delLogger)
		}

		// default option
		return listFiles(path, out)
	})
}
