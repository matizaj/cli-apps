package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

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
	pipeline := make([]step, 1)
	pipeline[0] = newStep(
		"go build",
		"go",
		"Go Build: SUCCESS",
		proj,
		[]string{"build", ".", "errors"},
	)
	for _, s := range pipeline {
		msg, err := s.execute()
		if err != nil {
			return err
		}

		_, err = fmt.Fprint(out, msg)
		if err != nil {
			return err
		}
	}
	return nil
}
