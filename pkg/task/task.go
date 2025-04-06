package task

import (
	"time"
)

// Priority represents the priority level of a task
type Priority int

const (
	Low Priority = iota
	Medium
	High
	Critical
)

// String returns the string representation of Priority
func (p Priority) String() string {
	switch p {
	case Low:
		return "LOW"
	case Medium:
		return "MEDIUM"
	case High:
		return "HIGH"
	case Critical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

// Status represents the current status of a task
type Status string

const (
	Pending   Status = "PENDING"
	Running   Status = "RUNNING"
	Completed Status = "COMPLETED"
	Failed    Status = "FAILED"
)

// Task represents a unit of work to be executed
type Task struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Command     string    `json:"command"`
	Priority    Priority  `json:"priority"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	StartedAt   time.Time `json:"started_at,omitempty"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	Error       string    `json:"error,omitempty"`
}

// TaskResult represents the result of a task execution
type TaskResult struct {
	TaskID      string    `json:"task_id"`
	Status      Status    `json:"status"`
	Output      string    `json:"output"`
	Error       string    `json:"error,omitempty"`
	CompletedAt time.Time `json:"completed_at"`
}

// TaskQueue is an interface for task queue operations
type TaskQueue interface {
	Enqueue(task Task) error
	Dequeue() (Task, error)
	Peek() (Task, error)
	Size() int
}

// TaskExecutor is an interface for task execution
type TaskExecutor interface {
	Execute(task Task) (TaskResult, error)
	Cancel(taskID string) error
}
