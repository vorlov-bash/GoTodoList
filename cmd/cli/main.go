package main

import (
	"flag"
	"fmt"
	"github.com/vorlov-bash/todolist/internal/buffers"
	"log"
	"os"
)

func printData(data []string) {
	for i, task := range data {
		fmt.Printf("%d. %s\n", i+1, task)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Print(
			"Available commands:\n" +
				"1. add - adds a task to a todo list\n" +
				"2. remove - removes task from a todo list\n" +
				"3. show - shows whole list\n",
		)
		os.Exit(1)
	}

	buffer, err := buffers.NewFileBuffer()
	if err != nil {
		log.Fatal(err)
	}

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	removeCmd := flag.NewFlagSet("remove", flag.ExitOnError)
	//showCmd := flag.NewFlagSet("show", flag.ExitOnError)

	addTask := addCmd.String("task", "", "Name of the task to add")
	removeTask := removeCmd.Int("task", 0, "Number of the task to remove")

	switch os.Args[1] {
	case "add":
		err := addCmd.Parse(os.Args[2:])

		if err != nil {
			log.Fatalf("Error parsing addCmd: %v", err)
		}

		if *addTask == "" {
			log.Fatal("Task cannot be empty")
		}

		data, err := buffer.Write(*addTask)
		if err != nil {
			log.Fatal(err)
		}
		printData(data)

	case "remove":
		err := removeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatalf("Error parsing removeCmd: %v", err)
		}

		if *removeTask == 0 {
			log.Fatal("Task cannot be empty")
		}

		data, err := buffer.Remove(*removeTask)
		if err != nil {
			log.Fatal(err)
		}
		printData(data)

	case "get":
		data, err := buffer.Get()

		if err != nil {
			log.Fatal(err)
		}
		printData(data)
	}
}
