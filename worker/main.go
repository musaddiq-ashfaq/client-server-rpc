package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"

	worker "worker/matrix_operation"
	client "client/matrix_request"
)

// WorkerService defines an RPC service for matrix computations.
type WorkerService struct {
	Address string
}

// Compute handles matrix operation requests from the Coordinator.
func (w *WorkerService) Compute(req client.MatrixRequest, res *client.MatrixResponse) error {
	
	result, err := worker.PerformMatrixOperation(req)
	if err != nil {
		return err
	}

	*res = result
	return nil
}

// registerWithCoordinator sends this worker's address to the Coordinator.
func registerWithCoordinator(workerAddr string) {
	coordinatorAddr := "localhost:50051" // Coordinator's RPC address

	client, err := rpc.Dial("tcp", coordinatorAddr)
	if err != nil {
		log.Printf("Failed to connect to coordinator: %v", err)
		time.Sleep(3 * time.Second)
		registerWithCoordinator(workerAddr)
		return
	}
	defer client.Close()

	var reply string
	err = client.Call("MatrixService.AddWorker", workerAddr, &reply)
	if err != nil {
		log.Printf("Failed to register with coordinator: %v", err)
		return
	}

	log.Printf("Successfully registered with coordinator at %s", coordinatorAddr)
}

func main() {
	workerPort := 50052 // Worker starts on this port
	workerAddr := fmt.Sprintf("localhost:%d", workerPort)

	// Register the RPC service
	workerService := &WorkerService{Address: workerAddr}
	err := rpc.Register(workerService)
	if err != nil {
		log.Fatalf("Error registering worker RPC: %v", err)
	}

	// Start listening for Coordinator requests
	listener, err := net.Listen("tcp", workerAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Worker is running on %s...", workerAddr)

	// Register itself with the Coordinator
	go registerWithCoordinator(workerAddr)

	// Handle incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Connection error: %v", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
