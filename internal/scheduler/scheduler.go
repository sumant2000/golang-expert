package scheduler

import (
	"container/heap"
	"sync"
	"time"

	"github.com/sumantkhapre/taskmaster/pkg/metrics"
	"github.com/sumantkhapre/taskmaster/pkg/task"
)

// PriorityQueue implements heap.Interface and holds Tasks
type PriorityQueue []*task.Task

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// Higher priority tasks come first
	if pq[i].Priority != pq[j].Priority {
		return pq[i].Priority > pq[j].Priority
	}
	// For tasks with same priority, older tasks come first
	return pq[i].CreatedAt.Before(pq[j].CreatedAt)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*task.Task)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// Scheduler manages task scheduling and distribution
type Scheduler struct {
	tasks    PriorityQueue
	mu       sync.Mutex
	workers  map[string]workerInfo
	workerMu sync.RWMutex
}

type workerInfo struct {
	lastHeartbeat time.Time
	status        string
}

func NewScheduler() *Scheduler {
	s := &Scheduler{
		tasks:   make(PriorityQueue, 0),
		workers: make(map[string]workerInfo),
	}
	heap.Init(&s.tasks)

	// Initialize metrics
	metrics.TaskQueueLength.Set(0)
	metrics.WorkerCount.Set(0)

	return s
}

func (s *Scheduler) AddTask(t *task.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}
	if t.Status == "" {
		t.Status = task.Pending
	}

	heap.Push(&s.tasks, t)

	// Update metrics
	metrics.TasksSubmitted.Inc()
	metrics.TaskQueueLength.Set(float64(s.tasks.Len()))
	metrics.TasksByPriority.WithLabelValues(t.Priority.String()).Inc()
}

func (s *Scheduler) GetNextTask() *task.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.tasks.Len() == 0 {
		return nil
	}

	t := heap.Pop(&s.tasks).(*task.Task)
	metrics.TaskQueueLength.Set(float64(s.tasks.Len()))
	return t
}

func (s *Scheduler) RegisterWorker(workerID string) {
	s.workerMu.Lock()
	defer s.workerMu.Unlock()

	s.workers[workerID] = workerInfo{
		lastHeartbeat: time.Now(),
		status:        "IDLE",
	}

	metrics.WorkerCount.Set(float64(len(s.workers)))
}

func (s *Scheduler) UpdateWorkerStatus(workerID, status string) {
	s.workerMu.Lock()
	defer s.workerMu.Unlock()

	if info, exists := s.workers[workerID]; exists {
		info.status = status
		info.lastHeartbeat = time.Now()
		s.workers[workerID] = info
	}
}

func (s *Scheduler) RemoveInactiveWorkers(timeout time.Duration) {
	s.workerMu.Lock()
	defer s.workerMu.Unlock()

	now := time.Now()
	for id, info := range s.workers {
		if now.Sub(info.lastHeartbeat) > timeout {
			delete(s.workers, id)
		}
	}

	metrics.WorkerCount.Set(float64(len(s.workers)))
}

func (s *Scheduler) GetWorkerCount() int {
	s.workerMu.RLock()
	defer s.workerMu.RUnlock()
	return len(s.workers)
}

func (s *Scheduler) GetQueueLength() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.tasks.Len()
}
