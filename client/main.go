package main

import (
	"fmt"
	"log"
	"net/rpc"
	mr "client/matrix_request"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:50051") // Connect to coordinator
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()

	req := mr.MatrixRequest{
		Operation: "add",
		MatrixA:   []int{1, 2, 3, 4},
		MatrixB:   []int{5, 6, 7, 8},
		RowsA:     2,
		ColsA:     2,
		RowsB:     2,
		ColsB:     2,
	}

	var res mr.MatrixResponse
	err = client.Call("MatrixService.Compute", req, &res) // Make an RPC call
	if err != nil {
		log.Fatalf("RPC error: %v", err)
	}
	fmt.Println("Response from coordinator:", res.Message)
	fmt.Println("Result:", res.Result)
}
