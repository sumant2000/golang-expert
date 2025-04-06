package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sumantkhapre/taskmaster/pkg/task"
)

var (
	serverAddr = flag.String("server", "http://localhost:8080", "TaskMaster server address")
	command    = flag.String("cmd", "", "Command to execute")
	priority   = flag.String("priority", "LOW", "Task priority (LOW, MEDIUM, HIGH, CRITICAL)")
	name       = flag.String("name", "", "Task name")
	desc       = flag.String("desc", "", "Task description")
)

func main() {
	flag.Parse()

	if *command == "" {
		fmt.Println("Error: command is required")
		flag.Usage()
		os.Exit(1)
	}

	if *name == "" {
		*name = fmt.Sprintf("task-%d", time.Now().Unix())
	}

	// Convert priority string to Priority type
	var taskPriority task.Priority
	switch *priority {
	case "LOW":
		taskPriority = task.Low
	case "MEDIUM":
		taskPriority = task.Medium
	case "HIGH":
		taskPriority = task.High
	case "CRITICAL":
		taskPriority = task.Critical
	default:
		fmt.Printf("Invalid priority: %s\n", *priority)
		os.Exit(1)
	}

	// Create task request
	taskReq := struct {
		Name        string        `json:"name"`
		Description string        `json:"description"`
		Command     string        `json:"command"`
		Priority    task.Priority `json:"priority"`
	}{
		Name:        *name,
		Description: *desc,
		Command:     *command,
		Priority:    taskPriority,
	}

	// Convert task to JSON
	jsonData, err := json.Marshal(taskReq)
	if err != nil {
		fmt.Printf("Error encoding task: %v\n", err)
		os.Exit(1)
	}

	// Submit task
	resp, err := http.Post(*serverAddr+"/tasks", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error submitting task: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		os.Exit(1)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error from server: %s\n", body)
		os.Exit(1)
	}

	// Parse response
	var result struct {
		TaskID string `json:"task_id"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Task submitted successfully!\nTask ID: %s\nStatus: %s\n", result.TaskID, result.Status)
}
