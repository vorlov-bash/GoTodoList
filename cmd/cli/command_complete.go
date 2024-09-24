package main

import (
	"flag"
	"github.com/vorlov-bash/todolist/pkg/cli"
	"github.com/vorlov-bash/todolist/pkg/tasks"
	"log"
	"os"
)

func CommandMarkAsComplete(buff tasks.Buffer) {
	// Add command
	completeCmd := flag.NewFlagSet("add", flag.ExitOnError)
	parsedTaskIdentifier := completeCmd.String("id", "", "Id of the task to remove")
	err := completeCmd.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}

	cli.FlagMustBeMandatory("id", *parsedTaskIdentifier)

	// Convert the string to an int
	taskIdentifier, err := cli.StringToInt(*parsedTaskIdentifier)

	if err != nil {
		log.Fatalf("Error parsing int: %v", err)
	}

	err = tasks.MarkAsDone(taskIdentifier, buff)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Task marked as completed: %v\n", taskIdentifier)
}
