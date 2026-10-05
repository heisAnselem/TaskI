package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"text/tabwriter"
	"time"
)

// Time computes the current time based on my format
func Time() string {
	layout := "Mon Jan _2 2006 3:04:05 PM MST"
	now := time.Now()
	return now.Format(layout)
}

// Task object - details regarding a task
type Task struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Status      Status `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type Status string

const (
	unknown    status = ""
	Todo       Status = "todo"
	InProgress Status = "in progress"
	Done       Status = "done"
)

// getNextId gets the next task id , if there are no tasks,the next id becomes 1.
// it keeps updating the max id in the tasks until no more tasks
func getNextId(tasks []Task) int {
	maxId := 0
	for _, task := range tasks {
		if task.Id > maxId {
			task.Id = maxId
		}
	}
	return maxId + 1
}

// getFilePath finds the absolute path for the JSON file in the current directory.
func getFilePath(filename string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, filename), nil
}

func readTask() ([]Task, error) {
	path, err := getFilePath("tasks.json")
	if err != nil {
		return nil, fmt.Errorf("failed to get file path: %w", err)
	}

	f, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// if file does not exist , creates empty task slice
			return []Task{}, nil
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	if len(f) == 0 {
		//file exists but no content
		return []Task{}, nil
	}
	var task []Task
	err = json.Unmarshal(f, &task)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return task, nil
}
func writeTask(task []Task) error {
	path, err := getFilePath("tasks.json")
	if err != nil {
		return fmt.Errorf("failed to get file path: %w", err)
	}
	data, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format JSON: %w", err)
	}
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

//CreateTask creates a task with description provided
func CreateTask(description string) {
	tasks, err := readTask()
	if err != nil {
		fmt.Printf("Error reading tasks : %w \n", err)
	}
	Id := getNextId(tasks)
	newTask := Task{
		Id:          Id,
		Description: description,
		Status:      Todo,
		CreatedAt:   Time(),
		UpdatedAt:   Time(),
	}
	tasks = append(tasks, newTask)
	err = writeTask(tasks)
	if err != nil {
		fmt.Printf("Error saving task %w", err)
	}
	fmt.Printf("Task created successfully (Id : %v ) \n", newTask.Id)
}

// UpdateTask updates a tasks description or status or both by id
func UpdateTask(Id int, description string, status Status) {
	tasks, err := readTask()
	if err != nil {
		fmt.Printf("Error reading tasks: %w\n", err)
	}
	if len(tasks) == 0 {
		// todo
		// show how to create tasks in print statement
		fmt.Printf("No tasks to Update\n")
	}
	if Id == 0 {
		fmt.Printf("Task id needs to be provided to update task\n")
	}
	if description == "" && status == unknown {
		fmt.Printf("At least description or status must be provided to update task \n")
	}
	for _, task := range tasks {
		if task.Id == Id {
			if description != "" {
				task.Description = description
			}
			if status != unknown {
				task.Status = status
			}
			task.UpdatedAt = Time()
			break
		}
	}
	err = writeTask(tasks)
	if err != nil {
		fmt.Printf("Failed to write task %w", err)
	}
	fmt.Println("Task updated successfully ")
	return nil
}

// DeleteTask deletes a task by its id
func DeleteTask(Id int) {
	tasks, err != readTask()
	if err != nil {
		fmt.Printf("Error reading tasks %w\n", err)
	}
	if Id == 0 {
		fmt.Printf("Task id needs to be provided to update task\n")
	}
	for i, task := range tasks {
		if task.Id == Id {
			// Delete elements from index i up to i+1 excluding i+1
			tasks = slices.Delete(tasks, i, i+1)
			fmt.Printf("Task by Id: %d Deleted successfully \n", Id)
			break
		}
	}
	fmt.Printf("No tasks with task id %d \n", Id)
	err = writeTask(tasks)
	if err != nil {
		fmt.Printf("Error Updating task memory base %w\n", err)
	}
}

// ListTasks lists tasks by their status , if status is not provided it Lists all tasks
func ListTasks(status Status) []Task {
	tasks, err := readTask()
	if err != nil {
		fmt.Printf("Error reading Tasks %w\n", err)
	}
	if len(tasks) == 0 {
		// todo
		// show how to create tasks in print statement
		fmt.Printf("No tasks to Update\n")
	}
	if status == unknown {
		return tasks
	}
	var taskList []Task
	for _, task := range tasks {
		if task.Status == status {
			taskList = append(taskList, task)
		}
	}
	return taskList
}

// FormatTask formats the tasks listed in a slice and prints to the terminal
func FormatTask(tasks []Task) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	fmt.Fprintln(w, "      TASK         ")
	fmt.Fprintln(w, "-------------------")
	// column headers
	fmt.Fprintln(w, "Id\tDescription\tStatus\tCreatedAt\tUpdatedAt\t")
	fmt.Fprintln(w, "--\t-----------\t------\t---------\t---------\t")

	//data rows
	for _, task := range tasks {
		s := fmt.Sprintf("%d\t%s\t%s\t%s\t%s\t", task.Id, task.Description, task.Status, task.CreatedAt, task.UpdatedAt)
		fmt.Fprintln(w, s)
	}
	w.Flush()
}
