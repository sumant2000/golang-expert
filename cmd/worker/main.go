package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/sumantkhapre/taskmaster/proto"
)

var (
	masterAddr = flag.String("master", "localhost:50051", "Master server address")
	workerId   = flag.String("worker-id", "", "Unique worker ID")
)

type worker struct {
	client pb.TaskServiceClient
	id     string
}

func newWorker(masterAddr, workerId string) (*worker, error) {
	conn, err := grpc.Dial(masterAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &worker{
		client: pb.NewTaskServiceClient(conn),
		id:     workerId,
	}, nil
}

func (w *worker) executeTask(task *pb.TaskRequest) (*pb.TaskResponse, error) {
	// Create a command from the task
	cmd := exec.Command("sh", "-c", task.Command)

	// Capture command output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &pb.TaskResponse{
			TaskId:  task.Id,
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.TaskResponse{
		TaskId:  task.Id,
		Success: true,
		Message: string(output),
	}, nil
}

func (w *worker) start() error {
	// TODO: Implement worker registration with master
	// TODO: Implement task polling and execution
	// For now, just keep the worker alive
	for {
		time.Sleep(time.Second)
	}

	return nil
}

func main() {
	flag.Parse()

	if *workerId == "" {
		hostname, err := os.Hostname()
		if err != nil {
			log.Fatalf("failed to get hostname: %v", err)
		}
		*workerId = hostname
	}

	w, err := newWorker(*masterAddr, *workerId)
	if err != nil {
		log.Fatalf("failed to create worker: %v", err)
	}

	log.Printf("Worker %s starting, connecting to master at %s", *workerId, *masterAddr)
	if err := w.start(); err != nil {
		log.Fatalf("worker failed: %v", err)
	}
}
