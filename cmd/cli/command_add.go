package main

import (
	"flag"
	"github.com/vorlov-bash/todolist/pkg/cli"
	"github.com/vorlov-bash/todolist/pkg/tasks"
	"log"
	"os"
	"time"
)

func CommandAdd(buff tasks.Buffer) {
	// Add command
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)

	taskName := addCmd.String("name", "", "Name of the task to add")
	taskDescription := addCmd.String("description", "", "Description of the task to add")
	taskDueDate := addCmd.String("due", "", "Due date of the task to add in format RFC3339")

	err := addCmd.Parse(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}

	cli.FlagMustBeMandatory("name", *taskName)
	cli.FlagMustBeMandatory("due", *taskDueDate)

	parsedDueDate, err := time.Parse(time.DateOnly, *taskDueDate)

	if err != nil {
		log.Fatalf("Error parsing date: %v", err)
	}

	taskOptions := tasks.TaskOptions{
		Name:        *taskName,
		Description: *taskDescription,
		DueDate:     parsedDueDate,
	}

	_, err = tasks.InsertTask(taskOptions, buff)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Task added: %v\n", taskOptions)
}
