package main

import (
	"crypto/tls"
	"log"
	"net"
	"net/rpc"
	"time"

	client "client/matrix_request"
	worker "worker/matrix_operation"
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

// Register with the Coordinator using TLS.
func registerWithCoordinator(workerAddr string) {
	coordinatorAddr := "localhost:50051" // Coordinator's RPC address

	config := &tls.Config{InsecureSkipVerify: true}
	client, err := tls.Dial("tcp", coordinatorAddr, config)
	if err != nil {
		log.Printf("Failed to connect to coordinator: %v", err)
		time.Sleep(3 * time.Second)
		registerWithCoordinator(workerAddr)
		return
	}
	defer client.Close()

	var reply string
	err = rpc.NewClient(client).Call("MatrixService.AddWorker", workerAddr, &reply)
	if err != nil {
		log.Printf("Failed to register with coordinator: %v", err)
		return
	}

	log.Printf("Successfully registered with coordinator at %s", coordinatorAddr)
}

func main() {
	// Start listening for Coordinator requests using TLS
	workerHost := "localhost"
	listener, err := net.Listen("tcp", workerHost+":0")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	
	workerAddr := listener.Addr().String()
	workerService := &WorkerService{Address: workerAddr}

	err = rpc.Register(workerService)
	if err != nil {
		log.Fatalf("Error registering worker RPC: %v", err)
	}

	log.Printf("Worker is running on %s...", workerAddr)

	// Register itself with the Coordinator
	go registerWithCoordinator(workerAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Connection error: %v", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
