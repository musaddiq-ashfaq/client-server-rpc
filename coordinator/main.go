package main

import (
	"log"
	"net"
	"net/rpc"
	ms "coordinator/matrix_service"
)

func main() {
	matrixService := new(ms.MatrixService)
	err := rpc.Register(matrixService) // Register the service
	if err != nil {
		log.Fatalf("Error registering RPC service: %v", err)
	}

	listener, err := net.Listen("tcp", ":50051") // Listen on port 50051
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("Coordinator is running on port 50051...")

	for {
		conn, err := listener.Accept() // Accept incoming connections
		if err != nil {
			log.Printf("Connection error: %v", err)
			continue
		}
		go rpc.ServeConn(conn) // Handle connection
	}
}
