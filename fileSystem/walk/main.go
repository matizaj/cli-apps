package main

import (
	"flag"
	"fmt"
	"os"
)

type config struct {
	// extension to filter out
	ext string

	// min file size
	size int64

	// list files
	list bool
}

func main() {
	root := flag.String("root", ".", "root directory")
	list := flag.Bool("list", false, "list files")
	ext := flag.String("ext", "", "file extension")
	size:=flag.Int64("size", 0, "minimum file size")
	

	flag.Parse()

	c:= config {
		ext: *ext,
		size: *size,
		list: *list,
	}

	if err := run(*root, os.Stdout, c); err!=nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(s string, file *os.File, c config) error {
	return nil
}