package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	pipeline := make([]executer, 4)
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
	pipeline[3] = newTimeoutStep(
		"git push",
		"git",
		"Git Push: SUCCESS",
		proj,
		[]string{"push", "origin", "master"},
		time.Second*3,
	)

	sig:=make(chan os.Signal, 1)
	errChan:= make(chan error)
	doneChan:=make(chan struct{})

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for _, s := range pipeline {
		msg, err := s.execute()
		if err != nil {
			errChan <- err
			return
		}

		_, err = fmt.Fprintln(out, msg)
		if err != nil {
			errChan<-err
			return
		}
	}
	close(doneChan)
	}()


	for {
		select {
			case err:= <-errChan:
				return err
			case rec := <-sig:
				signal.Stop(sig)
				return fmt.Errorf("%s: exiting %w", rec, ErrSignal)
			case <-doneChan:
				return nil
		}
	}
	
}
