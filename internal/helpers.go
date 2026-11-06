package internal

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/heisAnselem/TaskI/internal/task"
)

// Time computes the current time based on my format
func Time() string {
	layout := "Mon Jan _2 2006 3:04:05 PM MST"
	now := time.Now()
	return now.Format(layout)
}

// custom flag.typeVar function for my Status type
func StatusVar(fs *flag.FlagSet, p *task.Status, name string, value task.Status, usage string) {
	*p = value
	fs.Var(p, name, usage)
}

// getNextId gets the next task id , if there are no tasks,the next id becomes 1.
// it keeps updating the max id in the tasks until no more tasks
func getNextId(tasks []task.Task) int {
	maxId := 0
	for _, task := range tasks {
		if task.Id > maxId {
			maxId = task.Id
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

func readTask() ([]task.Task, error) {
	path, err := getFilePath("tasks.json")
	if err != nil {
		return nil, fmt.Errorf("failed to get file path: %w", err)
	}

	f, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// if file does not exist , creates empty task slice
			return []task.Task{}, nil
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	if len(f) == 0 {
		//file exists but no content
		return []task.Task{}, nil
	}
	task := []task.Task{}
	err = json.Unmarshal(f, &task)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return task, nil
}
func writeTask(task []task.Task) error {
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

// CreateTask creates a task with description provided
func CreateTask(description string) (task.Task, error) {
	if strings.TrimSpace(description) == "" {
		return task.Task{}, errors.New("task description is required")
	}
	tasks, err := readTask()
	if err != nil {
		return task.Task{}, fmt.Errorf("failed to read tasks: %w ", err)
	}
	Id := getNextId(tasks)
	newTask := task.Task{
		Id:          Id,
		Description: description,
		Status:      task.Todo,
		CreatedAt:   Time(),
		UpdatedAt:   Time(),
	}
	tasks = append(tasks, newTask)
	if err := writeTask(tasks); err != nil {
		return task.Task{}, fmt.Errorf("Error saving task %w", err)
	}
	return newTask, nil
}

// UpdateTask updates a tasks description or status or both by id,
// returns the updated task and error if any encountered
func UpdateTask(Id int, description string, status task.Status) (task.Task, error) {
	if Id <= 0 {
		return task.Task{}, errors.New("Task id needs to be provided to update task")
	}
	if strings.TrimSpace(description) == "" && status == task.Unknown {
		return task.Task{}, errors.New("At least description or status must be provided to update task")
	}
	tasks, err := readTask()
	if err != nil {
		return task.Task{}, fmt.Errorf("failed to read tasks: %w", err)
	}
	if len(tasks) == 0 {
		return task.Task{}, errors.New("No tasks available to update")
	}
	taskUpdated := task.Task{}
	found := false
	for i := range tasks {
		if tasks[i].Id == Id {
			found = true
			if strings.TrimSpace(description) != "" {
				tasks[i].Description = description
			}
			if status != task.Unknown {
				tasks[i].Status = status
			}
			tasks[i].UpdatedAt = Time()
			taskUpdated = tasks[i]
			break
		}
	}
	if !found {
		return task.Task{}, errors.New("Task not found ")
	}
	if err := writeTask(tasks); err != nil {
		return task.Task{}, fmt.Errorf("failed to write task %w", err)
	}
	return taskUpdated, nil
}

// DeleteTask deletes a task by its id
func DeleteTask(Id int) error {
	if Id <= 0 {
		return errors.New("task Id must be a positive integer")
	}
	tasks, err := readTask()
	if err != nil {
		return fmt.Errorf("failed to read tasks :%w", err)
	}
	for i, task := range tasks {
		if task.Id == Id {
			// Delete elements from index i up to i+1 excluding i+1
			tasks = slices.Delete(tasks, i, i+1)
			err = writeTask(tasks)
			if err != nil {
				return fmt.Errorf("Failed to write tasks : %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("tasks with Id %d not found", Id)
}

// ListTasks lists tasks by their status , if status is not provided it Lists all tasks
func ListTasks(status task.Status) ([]task.Task, error) {
	if status != task.Unknown && status != task.Todo && status != task.InProgress && status != task.Done {
		return []task.Task{}, fmt.Errorf("task Status must be either \"todo\",\"done\", or \"in progress\"")
	}
	tasks, err := readTask()
	if err != nil {
		return []task.Task{}, fmt.Errorf("failed to read tasks %w", err)
	}
	if len(tasks) == 0 || status == task.Unknown {
		return tasks, nil
	}
	taskList := []task.Task{}
	for _, task := range tasks {
		if task.Status == status {
			taskList = append(taskList, task)
		}
	}
	return taskList, nil
}

// FormatTask formats the tasks listed in a slice and prints to the terminal
func FormatTask(tasks []task.Task) {
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
