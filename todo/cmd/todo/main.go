package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"matizaj/cli-apps/todo"
	"os"
	"strings"
)

var todoFilename=".todo.json"

func main() {
	var add bool
	var completed int
	var list bool
	var del int
	var verbose bool

	if os.Getenv("TODO_FILENAME") != "" {
		todoFilename = os.Getenv("TODO_FILENAME")
	}

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
		"\n%s tool. Developed for learning\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2026\n")
	fmt.Fprintf(flag.CommandLine.Output(), "Usage information\n")
		flag.PrintDefaults()
	}


	flag.BoolVar(&add, "add",false, "task to include in ToDo list")
	flag.IntVar(&completed, "completed", 0, "task to mark as completed")
	flag.BoolVar(&list, "list", false, "list all tasks")
	flag.IntVar(&del, "del", -1, "delete selected item by index")
	flag.BoolVar(&verbose, "v", false, "display all propeties")

	flag.Parse()

	l:=&todo.List{}

	if err := l.Get(todoFilename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case list: 
		fmt.Print(l)
	case len(os.Args) ==1:
		for _, item := range *l {
			fmt.Println(item.Task)
		}
	case completed> 0: 
		if err := l.Complete(completed); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := l.Save(todoFilename); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case add:
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}		
		l.Add(t)

		if err := l.Save(todoFilename); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case del >= 0:
		if err := l.Delete(del); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		
		if err := l.Save(todoFilename); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case verbose: 
		txt, err := l.Verbose(todoFilename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(txt)
	default:
		fmt.Fprintln(os.Stderr, "invalid option")
		os.Exit(1)
	}
}


func getTask(r io.Reader, args ...string)(string, error) {
	if len(args)>0 {
		return strings.Join(args, " "), nil
	}

	s:=bufio.NewScanner(r)
	s.Scan()

	if err := s.Err(); err != nil {
		return "", err
	}

	if len(s.Text()) == 0 {
		return "", fmt.Errorf("ask cannot be blank")
	}

	return s.Text(), nil
}