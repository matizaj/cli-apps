package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	fmt.Println("reading from file..")
	flag.Parse()
	args := flag.Args()
	
	fmt.Printf("Args: %s\n", args)
	buffer := make([]byte, 4096)    

	for _, fname := range args {
		file, err := os.Open(fname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot open file: %v", err)
			os.Exit(1)
		}
		
		for {
			bytes, err := file.Read(buffer)
			if bytes >0 {
				fmt.Printf("%s ",string(buffer[:bytes]))
			}
			if err == io.EOF {
				break
			}
		}
	}
	
}
