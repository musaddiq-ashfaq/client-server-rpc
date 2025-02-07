package main

import (
	"log"
	"net"
	"net/rpc"

	worker "worker/matrix_operation" 
	client "client/matrix_request"   
)

// defines an RPC service.
type WorkerService struct{}

// Compute handles matrix operation requests from the Coordinator.
func (w *WorkerService) Compute(req client.MatrixRequest, res *client.MatrixResponse) error {
	
	result, err := worker.PerformMatrixOperation(req)
	if err != nil {
		return err
	}

	*res = result
	return nil
}

func main() {
	// Register the RPC service
	workerService := new(WorkerService)
	err := rpc.Register(workerService)
	if err != nil {
		log.Fatalf("Error registering worker RPC: %v", err)
	}

	// Start listening for Coordinator requests
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Worker is running on port 50052...")

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
