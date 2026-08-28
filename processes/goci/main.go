package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type executer interface {
	execute() (string, error)
}

func main() {
	fmt.Println("..::GOCI::..")

	proj := flag.String("p", "", "project directory path")
	flag.Parse()

	if err := run(*proj, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(proj string, out io.Writer) error {
	pipeline := make([]executer, 3)
	pipeline[0] = newStep(
		"go build",
		"go",
		"Go Build: SUCCESS",
		proj,
		[]string{"build", ".", "errors"},
	)
	pipeline[1] = newStep(
		"go test",
		"go",
		"Go Test: SUCCESS",
		proj,
		[]string{"test", "-v"},
	)
	pipeline[2] = newExceptionStep(
		"go fmt",
		"gofmt",
		"GoFmt: SUCCESS",
		proj,
		[]string{"-l", "."},
	)
	
	for _, s := range pipeline {
		msg, err := s.execute()
		if err != nil {
			return err
		}

		_, err = fmt.Fprintln(out, msg)
		if err != nil {
			return err
		}
	}
	return nil
}
