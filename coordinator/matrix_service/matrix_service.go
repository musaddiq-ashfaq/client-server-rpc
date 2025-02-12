package matrix_service

import (
	client "client/matrix_request"
	"errors"
	//"fmt"
	"log"
	"net/rpc"
	"sync"
)

// WorkerInfo holds information about a worker.
type WorkerInfo struct {
	Address string
	Load    int
}

// MatrixService defines an RPC service with a Compute method.
type MatrixService struct {
	workers []WorkerInfo
	mu      sync.Mutex
}

// NewMatrixService initializes a new MatrixService with worker addresses.
func NewMatrixService() *MatrixService {
	workers := []WorkerInfo{}
	/*for port := 50052; port <= 50060; port++ {
		address := fmt.Sprintf("localhost:%d", port)
		workers = append(workers, WorkerInfo{Address: address, Load: 0})
	}*/
	return &MatrixService{workers: workers}
}

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

// Compute handles matrix operations (for now, it just acknowledges the request).
func (m *MatrixService) Compute(req client.MatrixRequest, res *client.MatrixResponse) error {
	if req.Operation != "transpose" && (len(req.MatrixA) == 0 || len(req.MatrixB) == 0) {
		return errors.New("invalid matrices: both MatrixA and MatrixB are required for this operation")
	}
	if req.Operation == "transpose" && len(req.MatrixA) == 0 {
		return errors.New("invalid matrices: MatrixA is required for transpose operation")
	}

	log.Printf("Received request: Operation=%v, MatrixA=%v, MatrixB=%v\n", req.Operation, req.MatrixA, req.MatrixB)

	for {
		address, err := m.getLeastBusyWorker()
		if err != nil {
			log.Printf("Error getting least busy worker: %v", err)
			return errors.New("no available workers to handle the request")
		}

		client, err := rpc.Dial("tcp", address)
		if err != nil {
			log.Printf("Could not connect to worker at %s: %v", address, err)
			m.updateWorkerLoad(address, -1)
			continue
		}
		defer client.Close()

		m.updateWorkerLoad(address, 1)
		err = client.Call("WorkerService.Compute", req, res)
		m.updateWorkerLoad(address, -1)
		if err != nil {
			log.Printf("Error calling Compute on worker at %s: %v", address, err)
			continue
		}

		// Successfully sent the task to a worker
		log.Printf("Worker at %s returned: Result=%v, Rows=%d, Cols=%d, Message=%s", address, res.Result, res.Rows, res.Cols, res.Message)
		return nil
	}
}

