package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sumantkhapre/taskmaster/internal/scheduler"
	"github.com/sumantkhapre/taskmaster/pkg/task"
)

type Handler struct {
	scheduler *scheduler.Scheduler
}

func NewHandler(s *scheduler.Scheduler) *Handler {
	return &Handler{scheduler: s}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/tasks", h.submitTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", h.getTaskStatus).Methods("GET")
	r.HandleFunc("/tasks", h.listTasks).Methods("GET")
	r.HandleFunc("/workers", h.getWorkerStatus).Methods("GET")
}

type TaskRequest struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Command     string        `json:"command"`
	Priority    task.Priority `json:"priority"`
}

func (h *Handler) submitTask(w http.ResponseWriter, r *http.Request) {
	var req TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	t := &task.Task{
		ID:          generateTaskID(), // TODO: Implement ID generation
		Name:        req.Name,
		Description: req.Description,
		Command:     req.Command,
		Priority:    req.Priority,
		Status:      task.Pending,
		CreatedAt:   time.Now(),
	}

	h.scheduler.AddTask(t)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"task_id": t.ID,
		"status":  "submitted",
	})
}

func (h *Handler) getTaskStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	// TODO: Implement task status retrieval
	// For now, return a placeholder response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"task_id": taskID,
		"status":  "pending",
	})
}

func (h *Handler) listTasks(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement task listing
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"queue_length": h.scheduler.GetQueueLength(),
	})
}

func (h *Handler) getWorkerStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"worker_count": h.scheduler.GetWorkerCount(),
	})
}

// TODO: Implement proper task ID generation
func generateTaskID() string {
	return time.Now().Format("20060102150405.000")
}
