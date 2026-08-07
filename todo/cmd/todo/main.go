package main

import (
	"flag"
	"fmt"
	"matizaj/cli-apps/todo"
	"os"
)

var todoFilename=".todo.json"

func main() {
	var task string
	var completed int
	var list bool

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


	flag.StringVar(&task, "task","", "task to include in ToDo list")
	flag.IntVar(&completed, "completed", 0, "task to mark as completed")
	flag.BoolVar(&list, "list", false, "list all tasks")

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
	case task != "":
		l.Add(task)
		if err := l.Save(todoFilename); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	default:
		fmt.Fprintln(os.Stderr, "invalid option")
		os.Exit(1)
	}
}
