package main

import (
	ms "coordinator/matrix_service"
	"crypto/tls"
	"log"
	"net/rpc"
)

func main() {
	matrixService := ms.NewMatrixService(100)
	err := rpc.Register(matrixService) // Register the service
	if err != nil {
		log.Fatalf("Error registering RPC service: %v", err)
	}

	// Load TLS certificate
	cert, err := tls.LoadX509KeyPair("../certificates/cert.pem", "../certificates/key.pem")
	if err != nil {
		log.Fatalf("Failed to load TLS certificate: %v", err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	listener, err := tls.Listen("tcp", ":50051", config) // Listen with TLS
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("Coordinator is running on port 50051 with TLS...")

	for {
		conn, err := listener.Accept() // Accept incoming connections
		if err != nil {
			log.Printf("Connection error: %v", err)
			continue
		}
		go rpc.ServeConn(conn) // Handle connection
	}
}
