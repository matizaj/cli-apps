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
	return nil
}