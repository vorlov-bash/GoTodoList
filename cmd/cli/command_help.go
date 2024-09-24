package main

import (
	"fmt"
	"github.com/vorlov-bash/todolist/pkg/cli"
	"os"
)

func PrintHelp() {
	fmt.Print(
		"Usage: todolist <command> [arguments]\n\n" +
			"The commands are:\n\n" +
			"    add         Add a new task\n" +
			"    complete    Mark a task as completed\n" +
			"    remove      Remove a task\n" +
			"    show        Show all tasks\n" +
			"    help        Show this help\n\n" +
			"Use \"todolist help [command]\" for more information about a command.\n",
	)
}

func CommandHelp() {
	if len(os.Args) == 2 {
		PrintHelp()
		return
	}

	command := cli.Command(os.Args[2])

	if command == "" {
		fmt.Print("Command not found\n")
		return
	}

	if command == cli.CommandHelp {
		PrintHelp()
		return
	}

	switch command {
	case cli.CommandAdd:
		fmt.Print(
			"This command adds a new task\n\n"+
				"Usage: cli add [arguments]\n\n"+
				"The arguments are:\n\n"+
				"    -name         Name of the task to add\n"+
				"    -due          Due date of the task to add in date format\n"+
				"    -description  Description of the task to add (optional)\n\n",
			"Example: cli add -name \"Task 1\" -due \"2021-12-31\"\n",
		)
	case cli.CommandRemove:
		fmt.Print(
			"This command removes a task\n\n"+
				"Usage: cli remove [arguments]\n\n"+
				"The arguments are:\n\n"+
				"    -id         Id of the task to remove\n\n",
			"Example: cli remove -id 1\n",
		)
	case cli.CommandMarkAsCompleted:
		fmt.Print(
			"This command marks a task as completed\n\n"+
				"Usage: cli complete [arguments]\n\n"+
				"The arguments are:\n\n"+
				"    -id         Id of the task to mark as completed\n\n",
			"Example: cli complete -id 1\n",
		)
	case cli.CommandShow:
		fmt.Print(
			"This command shows all tasks\n\n"+
				"Usage: cli show\n\n",
			"Example: cli show\n",
		)
	default:
		fmt.Printf("Command %s not found\n", command)
	}
}
