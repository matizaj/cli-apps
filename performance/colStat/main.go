package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	fmt.Println("ColStat")

	op:=flag.String("op", "", "operation to perform")
	col:=flag.Int("col", 0, "CSV column number on which execute the operation")

	flag.Parse()

	if err := run(flag.Args(), *op, *col, os.Stdout); err!= nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

func run(filenames []string, op string, col int, out io.Writer) error {
	var opFunc statsFunc

	if len(filenames) ==0 {
		return ErrNoFiles
	}
	if col <1 {
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
	for _, fname := range filenames {
		f, err := os.Open(fname)
		if err != nil {
			return err
		}
		data, err := csv2float(f, col)
		if err != nil {
			return err
		}

		if err := f.Close(); err != nil {
			return err
		}

		consolidate = append(consolidate, data...)
	}
	_, err := fmt.Fprint(out, opFunc(consolidate))
	return err
}