package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"net/rpc"
	"os"
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

// Secure connection with TLS
func connectWithTLS() *rpc.Client {
	certPool := x509.NewCertPool()
	caCert, err := os.ReadFile("../certificates/ca-cert.pem") // Path to CA certificate
	if err != nil {
		log.Fatalf("Failed to read CA cert: %v", err)
	}
	certPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		RootCAs: certPool,
	}

	conn, err := tls.Dial("tcp", "localhost:50051", tlsConfig)
	if err != nil {
		log.Fatalf("Failed to connect with TLS: %v", err)
	}

	return rpc.NewClient(conn)
}

// registerWithCoordinator sends this worker's address to the Coordinator.
func registerWithCoordinator(workerAddr string) {
	for {
		client := connectWithTLS()
		if client == nil {
			time.Sleep(3 * time.Second)
			continue
		}

		var reply string
		err := client.Call("MatrixService.AddWorker", workerAddr, &reply)
		client.Close()
		if err != nil {
			log.Printf("Failed to register with coordinator: %v", err)
			time.Sleep(3 * time.Second)
			continue
		}

		log.Printf("Successfully registered with coordinator at %s", workerAddr)
		break
	}
}

func main() {

	// Start listening for Coordinator requests
	workerHost := "localhost"
	listener, err := net.Listen("tcp", workerHost+":0")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	workerAddr := listener.Addr().String()

	// Register the RPC service
	workerService := &WorkerService{Address: workerAddr}

	err = rpc.Register(workerService)
	if err != nil {
		log.Fatalf("Error registering worker RPC: %v", err)
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
