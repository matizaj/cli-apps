package main

import (
	//"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	file, _ := os.Open("test.txt")
	// if err != nil {
	// 	fmt.Printf("err: %v\n", err)
	// }
	// defer file.Close()

	// data, err := os.ReadFile(file.Name())
	// if err != nil {
	// 	fmt.Printf("cant read file %v\n", err)
	// }

	// fmt.Printf("%s", string(data))

	// scanner := bufio.NewScanner(file)
	// for scanner.Scan() {
	// 	fmt.Println(scanner.Text())
	// }
	buffer := make([]byte, 4096)
	counter:=0
	for{
		fmt.Println("line")
		readedChars, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		//fmt.Printf("counter %d", counter)
		counter++
		fmt.Println(readedChars)
		fmt.Println(string(buffer))
	}
	
	


}