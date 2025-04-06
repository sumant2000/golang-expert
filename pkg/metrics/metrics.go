package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// TasksSubmitted tracks the total number of submitted tasks
	TasksSubmitted = promauto.NewCounter(prometheus.CounterOpts{
		Name: "taskmaster_tasks_submitted_total",
		Help: "The total number of submitted tasks",
	})

	// TasksCompleted tracks the total number of completed tasks
	TasksCompleted = promauto.NewCounter(prometheus.CounterOpts{
		Name: "taskmaster_tasks_completed_total",
		Help: "The total number of completed tasks",
	})

	// TasksFailed tracks the total number of failed tasks
	TasksFailed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "taskmaster_tasks_failed_total",
		Help: "The total number of failed tasks",
	})

	// TaskQueueLength tracks the current length of the task queue
	TaskQueueLength = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "taskmaster_task_queue_length",
		Help: "The current length of the task queue",
	})

	// WorkerCount tracks the number of connected workers
	WorkerCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "taskmaster_worker_count",
		Help: "The number of connected workers",
	})

	// TaskExecutionTime tracks the time taken to execute tasks
	TaskExecutionTime = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "taskmaster_task_execution_seconds",
		Help:    "Time taken to execute tasks",
		Buckets: prometheus.ExponentialBuckets(0.1, 2.0, 10),
	})

	// TasksByPriority tracks tasks by priority level
	TasksByPriority = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "taskmaster_tasks_by_priority_total",
			Help: "Number of tasks by priority level",
		},
		[]string{"priority"},
	)
)
