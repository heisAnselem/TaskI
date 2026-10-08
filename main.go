package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/heisAnselem/TaskI/internal"
)

const (
	Unknown    = internal.Unknown
	Todo       = internal.Todo
	InProgress = internal.InProgress
	Done       = internal.Done
)

type Status = internal.Status
type Task = internal.Task

func main() {

	// checks if user didn't pass a subcommand
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Independent flag sets for each subcommand
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)

	task := Task{}
	// Binding Create command flags
	createCmd.StringVar(&task.Description, "description", "", "Task Description")

	// Binding list command flags
	internal.StatusVar(listCmd, &task.Status, "status", Unknown, "Task Status ")

	// Binding Update command flags
	updateCmd.IntVar(&task.Id, "id", 0, "Task id ")
	updateCmd.StringVar(&task.Description, "description", "", "Task Description")
	internal.StatusVar(updateCmd, &task.Status, "status", Unknown, "Task Status")

	//Binding Delete command flags
	deleteCmd.IntVar(&task.Id, "id", 0, "Task id ")

	// Parse the subcommand specific flags (skip program name and subcommand name)
	subcommand := os.Args[1]
	subArgs := os.Args[2:]

	switch subcommand {
	case "create":
		createCmd.Parse(subArgs)
		fmt.Printf("Action : CREATE -> Task\n")
		internal.CreateTask(task.Description)
	case "list":
		listCmd.Parse(subArgs)
		if task.Status == Unknown {
			fmt.Printf("Action : LIST -> All Tasks\n")
		} else {
			fmt.Printf("Action : LIST -> All Tasks by Status \n")
		}
		tasks := internal.ListTasks(task.Status)
		internal.FormatTask(tasks)
	case "update":
		updateCmd.Parse(subArgs)
		fmt.Printf("Action : UPDATE Task\n")
		internal.UpdateTask(task.Id, task.Description, task.Status)

	case "delete":
		deleteCmd.Parse(subArgs)
		fmt.Printf("Action : DELETE Task\n")
		internal.DeleteTask(task.Id)

	default:
		fmt.Printf("Unknown subcommand: %q\n\n", subcommand)
		printUsage()
		os.Exit(1)
	}

}

func printUsage() {
	fmt.Println("TaskI - A lightweight CLI task tracker")
	fmt.Println("\nUsage:")
	fmt.Println("  TaskI create -description \"Buy coffee   \"")
	fmt.Println("  TaskI list                           						    (List all tasks)")
	fmt.Println("  TaskI list -status \"done\"          						    (List all done tasks)")
	fmt.Println("  TaskI list -status \"in progress\"   						    (List all pending tasks)")
	fmt.Println("  TaskI update -id 1 -status \"in progress\"            		    (Update a task's status by task Id )")
	fmt.Println("  TaskI update -id 1 -description \"in progress\"                  (Update a task's description by task Id )")
	fmt.Println("  TaskI update -id 1 -status \"in progress\" -description \"Hello\"(Update a task's status or description by task Id )")
	fmt.Println("  TaskI delete -id 1            									(Delete a task)")
}
