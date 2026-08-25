package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
)

func main() {
	fmt.Println("ColStat")

	op := flag.String("op", "", "operation to perform")
	col := flag.Int("col", 0, "CSV column number on which execute the operation")

	flag.Parse()

	if err := run(flag.Args(), *op, *col, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func run(filenames []string, op string, col int, out io.Writer) error {
	resChan := make(chan []float64)
	errChan := make(chan error)
	doneChan := make(chan struct{})
	filesChan := make(chan string)

	var opFunc statsFunc
	var wg sync.WaitGroup

	if len(filenames) == 0 {
		return ErrNoFiles
	}
	if col < 1 {
		return ErrInvalidColumn
	}

	switch op {
	case "s":
		opFunc = sum

	case "a":
		opFunc = avg
	default:
		return ErrInvalidOperation
	}

	consolidate := make([]float64, 0)

	go func() {
		defer close(filesChan)
		for _, fname := range filenames {
			filesChan <- fname
		}
	}()
	
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for fname := range filesChan {
				f, err := os.Open(fname)
				if err != nil {
					errChan <- err
					return
				}
				data, err := csv2float(f, col)
				if err != nil {
					errChan <- err
					return
				}

				if err := f.Close(); err != nil {
					errChan <- err
					return
				}
				resChan <- data
			}
		}()
	}

	go func() {
		wg.Wait()
		close(doneChan)
	}()

	for {
		select {
		case err := <-errChan:
			return err
		case data := <-resChan:
			consolidate = append(consolidate, data...)
		case <-doneChan:
			_, err := fmt.Fprint(out, opFunc(consolidate))
			return err
		}
	}
}
