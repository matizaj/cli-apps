package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func filterOut(path, ext string, size int64, info os.FileInfo) bool {
	if info.IsDir() || info.Size() < size {
		return true
	}
	if ext != "" && filepath.Ext(path) != ext {
		return true
	}
	return false
}

func listFiles(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
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

	relDir, err := filepath.Rel(root, filepath.Dir(path))
	if err != nil {
		return err
	}

	dest:=fmt.Sprintf("%s.gz", filepath.Base(path))
	targetPath:= filepath.Join(destDir, relDir, dest)

	if err:=os.MkdirAll(filepath.Dir(targetPath), 0755); err!= nil {
		return err
	}

	return nil
}
