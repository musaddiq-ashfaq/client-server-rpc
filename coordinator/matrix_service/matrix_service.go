package matrix_service

import (
	mr "client/matrix_request"
	"errors"
	"log"
	"net/rpc"
	"sync"
)

// WorkerInfo holds information about a worker.
type WorkerInfo struct {
	Address string
	Load    int
}

// Task represents a matrix operation request along with a response channel.
type Task struct {
	Req  mr.MatrixRequest
	Resp chan mr.MatrixResponse
}

// MatrixService defines an RPC service which receives client requests and delegates them to workers.
type MatrixService struct {
	workers   []WorkerInfo
	mu        sync.Mutex
	taskQueue chan Task
}

// NewMatrixService initializes a new MatrixService with worker addresses and starts processing tasks.
func NewMatrixService(queueSize int) *MatrixService {
	m := &MatrixService{
		workers:   []WorkerInfo{},
		taskQueue: make(chan Task, queueSize),
	}

	// Start a worker goroutine to process tasks
	go m.processTasks()

	return m
}

// processTasks continuously consumes tasks from the queue and assigns them to the least busy worker.
func (m *MatrixService) processTasks() {
	for task := range m.taskQueue {
		for {
			address, err := m.getLeastBusyWorker()
			if err != nil {
				log.Printf("Error getting least busy worker: %v", err)
				task.Resp <- mr.MatrixResponse{Message: "No available workers"}
				break
			}

			client, err := rpc.Dial("tcp", address)
			if err != nil {
				log.Printf("Could not connect to worker at %s: %v", address, err)
				m.updateWorkerLoad(address, -1)
				continue
			}
			defer client.Close()

			var res mr.MatrixResponse
			err = client.Call("WorkerService.Compute", task.Req, &res)
			m.updateWorkerLoad(address, -1)
			if err != nil {
				log.Printf("Error calling Compute on worker at %s: %v", address, err)
				continue
			}

			// Send the response back to the original caller
			task.Resp <- res
			break
		}
	}
}

// Serve enqueues the request into the task queue for processing.
func (m *MatrixService) Serve(req mr.MatrixRequest, res *mr.MatrixResponse) error {
	if req.Operation != "transpose" && (len(req.MatrixA) == 0 || len(req.MatrixB) == 0) {
		return errors.New("invalid matrices: both MatrixA and MatrixB are required for this operation")
	}
	if req.Operation == "transpose" && len(req.MatrixA) == 0 {
		return errors.New("invalid matrices: MatrixA is required for transpose operation")
	}

	log.Printf("Received request: Operation=%v, MatrixA=%v, MatrixB=%v\n", req.Operation, req.MatrixA, req.MatrixB)

	// Create a response channel and enqueue the task
	respChan := make(chan mr.MatrixResponse, 1)
	m.taskQueue <- Task{Req: req, Resp: respChan}

	// Wait for response
	*res = <-respChan
	return nil
}

// AddWorker allows a Worker to register itself with the Coordinator.
func (m *MatrixService) AddWorker(workerAddr string, reply *string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if the worker is already registered
	for _, worker := range m.workers {
		if worker.Address == workerAddr {
			log.Printf("Worker at %s is already registered.", workerAddr)
			*reply = "Already registered"
			return nil
		}
	}

	// Add the new Worker
	m.workers = append(m.workers, WorkerInfo{Address: workerAddr, Load: 0})
	log.Printf("New Worker registered: %s", workerAddr)
	*reply = "Worker added successfully"
	return nil
}

// getLeastBusyWorker selects the worker with the lowest load.
func (m *MatrixService) getLeastBusyWorker() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.workers) == 0 {
		return "", errors.New("no available workers")
	}

	var leastBusyWorker *WorkerInfo
	for i := range m.workers {
		address := m.workers[i].Address

		// Check if the worker is reachable
		client, err := rpc.Dial("tcp", address)
		if err != nil {
			log.Printf("Skipping worker at %s: Connection failed", address)
			continue
		}
		client.Close()

		// Select the least busy worker
		if leastBusyWorker == nil || m.workers[i].Load < leastBusyWorker.Load {
			leastBusyWorker = &m.workers[i]
		}
	}

	if leastBusyWorker == nil {
		return "", errors.New("no available workers")
	}

	return leastBusyWorker.Address, nil
}

// updateWorkerLoad updates the load of a worker.
func (m *MatrixService) updateWorkerLoad(address string, delta int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, worker := range m.workers {
		if worker.Address == address {
			m.workers[i].Load += delta
			break
		}
	}
}