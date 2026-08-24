package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
)

func main() {
	fmt.Println("reading from file..")
	var wg sync.WaitGroup
	flag.Parse()
	args := flag.Args()
	
	fmt.Printf("Args: %s\n", args)
	resultChan := make(chan string)
	doneChan := make(chan struct{})

	
	

	for _, fname := range args {
		wg.Add(1)
		go func(fname string) {
			buffer := make([]byte, 4096)
			defer wg.Done()
			file, err := os.Open(fname)
			

			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot open file: %v", err)
				os.Exit(1)
			}
			defer file.Close()
			for {
				bytes, err := file.Read(buffer)
				if bytes >0 {
					resultChan <- string(buffer[:bytes])
				}
				if err == io.EOF {
					break
				}
			}
			
		}(fname)
	}
	
	go func() {
		wg.Wait()
		close(doneChan)
	}()

	for {
		select {
			case data:= <-resultChan: 
				fmt.Println(data)
			case <-doneChan:
				return
		}
	}
	
}
