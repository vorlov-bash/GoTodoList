package main

import (
	"flag"
	"github.com/vorlov-bash/todolist/pkg/cli"
	"github.com/vorlov-bash/todolist/pkg/tasks"
	"log"
	"os"
)

func CommandRemove(buff tasks.Buffer) {
	// Add command
	removeCmd := flag.NewFlagSet("add", flag.ExitOnError)
	parsedTaskIdentifier := removeCmd.String("id", "", "Id of the task to remove")
	err := removeCmd.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}

	cli.FlagMustBeMandatory("id", *parsedTaskIdentifier)

	// Convert the string to an int
	taskIdentifier, err := cli.StringToInt(*parsedTaskIdentifier)

	if err != nil {
		log.Fatalf("Error parsing int: %v", err)
	}

	err = tasks.DeleteTaskById(taskIdentifier, buff)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Task removed: %v\n", taskIdentifier)
}
