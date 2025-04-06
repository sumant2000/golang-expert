# TaskMaster - Distributed Task Scheduler

TaskMaster is a high-performance distributed task scheduling and execution system built in Go. It demonstrates Go's unique strengths in concurrent programming, network communication, and system design.

## Why Go?

This project specifically leverages Go's unique features:

1. **Concurrency**: Uses goroutines and channels for efficient task scheduling and execution
2. **Performance**: Go's compiled nature and efficient runtime make it ideal for high-throughput systems
3. **Network Programming**: Built-in support for gRPC and WebSocket makes distributed systems easy
4. **Standard Library**: Rich standard library reduces external dependencies
5. **Cross-Platform**: Easy deployment across different operating systems

## Features

- Distributed task scheduling and execution
- Real-time task status updates via WebSocket
- Load balancing across worker nodes
- Task prioritization and scheduling
- Health monitoring and metrics
- REST API for task management
- gRPC for internal communication

## Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Master     │     │  Worker 1   │     │  Worker N   │
│  (Scheduler)│◄───►│  (Executor) │◄───►│  (Executor) │
└─────────────┘     └─────────────┘     └─────────────┘
       ▲
       │
       ▼
┌─────────────┐
│  API Server │
└─────────────┘
```

## Building and Running

### Prerequisites

- Go 1.21 or later
- Protocol Buffers compiler (protoc)
- Make (optional)

### Building

1. Clone the repository:
   ```bash
   git clone https://github.com/sumantkhapre/taskmaster.git
   cd taskmaster
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build all components:
   ```bash
   go build ./...
   ```

### Running the System

1. Start the master node:
   ```bash
   go run cmd/master/main.go
   ```

2. Start one or more worker nodes:
   ```bash
   # Start worker 1
   go run cmd/worker/main.go --worker-id=worker1

   # Start worker 2
   go run cmd/worker/main.go --worker-id=worker2
   ```

3. Submit a task using the client:
   ```bash
   # Submit a simple task
   go run cmd/client/main.go --cmd="examples/task_example.sh" --priority=HIGH

   # Submit with custom name and description
   go run cmd/client/main.go \
     --cmd="examples/task_example.sh" \
     --priority=HIGH \
     --name="important-task" \
     --desc="An important task that needs to be executed"
   ```

### Monitoring

The system exposes metrics in Prometheus format at `http://localhost:8080/metrics`. You can use Prometheus and Grafana to visualize:

- Task queue length
- Number of active workers
- Task execution times
- Success/failure rates
- Tasks by priority

## API Endpoints

- `POST /tasks` - Submit a new task
- `GET /tasks/{id}` - Get task status
- `GET /tasks` - List all tasks
- `GET /workers` - Get worker status
- `GET /metrics` - Prometheus metrics

## Example Task

The repository includes an example task script (`examples/task_example.sh`) that demonstrates the system's capabilities:

```bash
# Run the example task directly
./examples/task_example.sh

# Submit it to TaskMaster
go run cmd/client/main.go --cmd="./examples/task_example.sh"
```

## Project Structure

```
taskmaster/
├── cmd/                # Main applications
│   ├── master/        # Master node
│   ├── worker/        # Worker nodes
│   └── client/        # Command-line client
├── internal/          # Private application code
│   ├── scheduler/     # Task scheduling logic
│   ├── executor/      # Task execution logic
│   ├── api/          # REST API handlers
│   └── grpc/         # gRPC service definitions
├── pkg/              # Public library code
│   ├── task/         # Task definitions
│   └── metrics/      # Monitoring and metrics
├── proto/            # Protocol buffer definitions
└── examples/         # Example tasks and scripts
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT
