package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Time computes the current time based on my format
func Time() string {
	time.Layout = "Mon Jan _2 2006 3:04:05 PM MST"
	now := time.Now()
	return now.Format(time.Layout)
}

// Task object - details regarding a task
type Task struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type Status string

const (
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
	err := json.Unmarshal(f, &task)
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
		fmt.Errorf("failed to format JSON: %w", err)
	}
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
}

func CreateTask(description string) {
	tasks, err := readTask()
	if err != nil {
		fmt.Errorf("failed to read tasks: %w", err)
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
		fmt.Errorf("Error saving task %w", err)
	}
	fmt.Printf("Task created successfully with Id : %v", newTask.Id)
}
func UpdateTask() {}

func DeleteTask() {}

func GetTaskById(id int) Task {}

func GetTaskByStatus(status Status) Task {}

func ListTasks() []Task {}
