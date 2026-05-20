package main

import (
	"fmt"
	"os"

	"github.com/Marif226/task-cli/internal/command"
	"github.com/Marif226/task-cli/internal/storage"
)

const (
	tasksDirPath  = "tasks"
	tasksFilePath = "tasks/tasks.json"
)

// Command describes a runnable CLI sub-command.
type Command interface {
	Name() string
	Init(args []string) error
	Run() error
}

func main() {
	err := os.MkdirAll(tasksDirPath, 0755)
	if err != nil {
		fmt.Printf("Error creating tasks directory: %s\n", err)
		os.Exit(1)
	}

	// Check if tasks file exists, if not create it
	_, err = os.Stat(tasksFilePath)
	if os.IsNotExist(err) {
		err := os.WriteFile(tasksFilePath, []byte(`{"last_id":0,"tasks":[]}`), 0644)
		if err != nil {
			fmt.Printf("Error creating %s: %s\n", tasksFilePath, err)
			os.Exit(1)
		}
	} else if err != nil {
		fmt.Printf("Error checking %s: %s\n", tasksFilePath, err)
		os.Exit(1)
	}

	taskStorage := storage.NewTaskStorage(tasksFilePath)

	cmdArg := os.Args[1:]

	if len(cmdArg) < 1 {
		fmt.Println("You must pass a sub-command")
		os.Exit(1)
	}

	cmds := []Command{
		command.NewAddCommand(taskStorage),
		command.NewDeleteCommand(taskStorage),
		command.NewListCommand(taskStorage),
		command.NewMarkCommand(taskStorage),
	}

	for _, cmd := range cmds {
		if cmd.Name() == cmdArg[0] {
			err := cmd.Init(cmdArg[1:])
			if err != nil {
				fmt.Printf("Error initializing command: %s\n", err)
				os.Exit(1)
			}

			err = cmd.Run()
			if err != nil {
				fmt.Printf("Error running command: %s\n", err)
				os.Exit(1)
			}

			return
		}
	}

	fmt.Printf("Unknown sub-command: %s\n", cmdArg[0])
	os.Exit(1)
}
