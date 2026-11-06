package task

import (
	"fmt"
	"strings"
)

// Task object - details regarding a task
type Task struct {
	Id          int    `json:"Id"`
	Description string `json:"Description"`
	Status      Status `json:"Status"`
	CreatedAt   string `json:"CreatedAt"`
	UpdatedAt   string `json:"UpdatedAt"`
}

type Status string

const (
	Unknown    Status = ""
	Todo       Status = "todo"
	InProgress Status = "in progress"
	Done       Status = "done"
)

// implementing the flag.value interface methods String(),Set()

func (s *Status) String() string {
	return string(*s)
}
func (s *Status) Set(value string) error {
	v := Status(strings.ToLower(value))
	switch v {
	case Unknown, Todo, InProgress, Done:
		*s = v
		return nil
	default:
		return fmt.Errorf("Status value must be either todo,in progress or done ")
	}
}
