package main

import (
	"github.com/vorlov-bash/todolist/pkg/cli"
	"github.com/vorlov-bash/todolist/pkg/tasks"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		PrintHelp()
		return
	}

	buffer, err := tasks.NewSqlite3Buffer("tmp/tasks.db")
	if err != nil {
		log.Fatal(err)
	}

	command := cli.Command(os.Args[1])

	switch command {
	case cli.CommandAdd:
		CommandAdd(buffer)
	case cli.CommandRemove:
		CommandRemove(buffer)
	case cli.CommandMarkAsCompleted:
		CommandMarkAsComplete(buffer)
	case cli.CommandShow:
		CommandShow(buffer)
	case cli.CommandHelp:
		CommandHelp()
	}
}
