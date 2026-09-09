package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	fmt.Println("ToDo Server")
	host := flag.String("h", "localhost", "Server host")
	port := flag.Int("p", 8080, "Server port")
	todoFile := flag.String("f", "todoSrv.json", "todo JSON file")

	flag.Parse()

	srv:=&http.Server{
		Addr: fmt.Sprintf("%s:%d", *host, *port),
		Handler: newMux(*todoFile),
		ReadTimeout: 10*time.Second,
		WriteTimeout: 10*time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
