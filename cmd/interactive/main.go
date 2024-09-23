package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/vorlov-bash/todolist/buffers"
	"log"
	"os"
	"strconv"
	"strings"
)

type Command string

const (
	commandAdd    Command = "add"
	commandRemove Command = "remove"
	commandShow   Command = "show"
	commandHelp   Command = "help"
	commandExit   Command = "exit"
)

func printData(data []string) {
	for i, task := range data {
		fmt.Printf("%d. %s\n", i+1, task)
	}
}

func parseCommand(scommand string) (Command, string, error) {
	splitSet := strings.Split(scommand, " ")
	command := Command(strings.Trim(splitSet[0], "\n"))

	switch command {
	case commandAdd:
		if len(splitSet) < 2 {
			return "", "", errors.New("task cannot be empty")
		}

		content := strings.Trim(strings.Join(splitSet[1:], " "), "\"")
		return commandAdd, content, nil
	case commandRemove:
		if len(splitSet) < 2 {
			return "", "", fmt.Errorf("task cannot be empty")
		}

		content := strings.Trim(strings.Join(splitSet[1:], " "), "\"")
		return commandRemove, content, nil
	case commandShow:
		return commandShow, "", nil
	case commandHelp:
		return commandHelp, "", nil
	case commandExit:
		return commandExit, "", nil
	default:
		return "", "", fmt.Errorf("unknown command: %s", command)
	}
}

func handleCommand(buffer buffers.Buffer, command Command, content string) error {
	switch command {
	case commandAdd:
		handleAdd(buffer, content)
	case commandRemove:
		number, err := strconv.Atoi(strings.Trim(content, "\n"))
		if err != nil {
			return fmt.Errorf("cannot convert content to number: %w", err)
		}
		handleRemove(buffer, number)
	case commandShow:
		handleShow(buffer)
	case commandHelp:
		handleHelp()
	case commandExit:
		handleExit()
	}

	return nil
}

func handleAdd(buffer buffers.Buffer, task string) {
	data, err := buffer.Write(task)
	if err != nil {
		log.Fatal(err)
	}
	printData(data)
}

func handleRemove(buffer buffers.Buffer, number int) {
	data, err := buffer.Remove(number)
	if err != nil {
		log.Fatal(err)
	}
	printData(data)
}

func handleShow(buffer buffers.Buffer) {
	data, err := buffer.Get()
	if err != nil {
		log.Fatal(err)
	}
	printData(data)
}

func handleHelp() {
	fmt.Print(
		"Available commands:\n" +
			"1. add <task> - adds a task to a todo list\n" +
			"2. remove <number> - removes task from a todo list\n" +
			"3. show - shows whole list\n" +
			"4. help - prints this message\n" +
			"5. exit - exits the program\n",
	)
}

func handleExit() {
	os.Exit(0)
}

func main() {
	buffer, err := buffers.NewSqlite3Buffer()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Hello! Welcome to the interactive todo list!\n" +
		"Type 'help' to see the list of available commands\n" +
		"Type 'exit' to exit the program\n" +
		"Let's get started!\n")

	for true {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		input_line, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			continue
		}

		command, content, err := parseCommand(input_line)
		if err != nil {
			if strings.Contains(err.Error(), "unknown command") {
				log.Println("Unknown command. Type 'help' to see the list of available commands")
			}
			continue
		}

		err = handleCommand(buffer, command, content)
		if err != nil {
			log.Println(err)
			continue
		}
	}

}
