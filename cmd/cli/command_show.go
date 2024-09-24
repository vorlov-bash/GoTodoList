package main

import (
	"github.com/vorlov-bash/todolist/pkg/cli"
	"github.com/vorlov-bash/todolist/pkg/tasks"
	"log"
)

func CommandShow(buff tasks.Buffer) {
	// Add command
	data, err := tasks.GetAllTasks(buff)

	if err != nil {
		log.Fatal(err)
	}

	cli.DisplayPrettyStruct(data)
}
