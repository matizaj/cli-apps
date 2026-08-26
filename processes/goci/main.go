package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("..::GOCI::..")

	proj := flag.String("p","", "project directory path")
	flag.Parse()

	if err := run(*proj, os.Stdout);err!=nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}


func run(proj string, out io.Writer) error {
	if proj == "" {
		return fmt.Errorf("project directory is required %w", ErrValidation)
	}
	args := []string{"build", ".", "error"}

	cmd := exec.Command("go", args...)
	cmd.Dir = proj

	if err:= cmd.Run(); err!= nil {
		return fmt.Errorf("'go build' failed: %s", err)
	}

	_, err := fmt.Fprintln(out, "Go Build: SUCCESS")

	return err
}