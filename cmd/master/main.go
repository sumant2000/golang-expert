package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	"github.com/sumantkhapre/taskmaster/internal/api"
	"github.com/sumantkhapre/taskmaster/internal/scheduler"
	pb "github.com/sumantkhapre/taskmaster/proto"
)

var (
	grpcPort = flag.String("grpc-port", ":50051", "gRPC server port")
	httpPort = flag.String("http-port", ":8080", "HTTP server port")
)

type masterServer struct {
	pb.UnimplementedTaskServiceServer
	scheduler *scheduler.Scheduler
}

func (s *masterServer) SubmitTask(ctx context.Context, req *pb.TaskRequest) (*pb.TaskResponse, error) {
	// TODO: Convert proto task to internal task and submit to scheduler
	return &pb.TaskResponse{
		TaskId:  req.Id,
		Success: true,
		Message: "Task submitted successfully",
	}, nil
}

func (s *masterServer) GetTaskStatus(ctx context.Context, req *pb.TaskStatusRequest) (*pb.TaskStatusResponse, error) {
	// TODO: Implement task status retrieval from scheduler
	return &pb.TaskStatusResponse{
		TaskId: req.TaskId,
		Status: "PENDING",
	}, nil
}

func (s *masterServer) StreamTaskUpdates(req *pb.TaskStatusRequest, stream pb.TaskService_StreamTaskUpdatesServer) error {
	// TODO: Implement task update streaming
	return nil
}

func main() {
	flag.Parse()

	// Create scheduler
	s := scheduler.NewScheduler()

	// Create gRPC server
	lis, err := net.Listen("tcp", *grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterTaskServiceServer(grpcServer, &masterServer{scheduler: s})

	go func() {
		log.Printf("gRPC server listening on %s", *grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Create HTTP server for API and metrics
	r := mux.NewRouter()

	// Register API routes
	apiHandler := api.NewHandler(s)
	apiHandler.RegisterRoutes(r)

	// Register metrics endpoint
	r.Handle("/metrics", promhttp.Handler())

	// Start worker cleanup routine
	go func() {
		for {
			s.RemoveInactiveWorkers(5 * time.Minute)
			time.Sleep(time.Minute)
		}
	}()

	// Start HTTP server
	httpServer := &http.Server{
		Addr:    *httpPort,
		Handler: r,
	}

	log.Printf("HTTP server listening on %s", *httpPort)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("failed to serve HTTP: %v", err)
	}
}
