package matrix_service

import (
	client "client/matrix_request"
	"errors"
	"fmt"
	"log"
	"net/rpc"
)

// MatrixService defines an RPC service with a Compute method.
type MatrixService struct{}

// Compute handles matrix operations (for now, it just acknowledges the request).
func (m *MatrixService) Compute(req client.MatrixRequest, res *client.MatrixResponse) error {
	if len(req.MatrixA) == 0 || len(req.MatrixB) == 0 {
		return errors.New("invalid matrices")
	}
	log.Printf("Received request: Operation=%v, MatrixA=%v, MatrixB=%v\n", req.Operation, req.MatrixA, req.MatrixB)

	// Try to connect to workers from port 50052 to 50060
	for port := 50052; port <= 50060; port++ {
		address := fmt.Sprintf("localhost:%d", port)
		client, err := rpc.Dial("tcp", address)
		if err != nil {
			log.Printf("Could not connect to worker at %s: %v", address, err)
			continue
		}
		defer client.Close()

		err = client.Call("WorkerService.Compute", req, res)
		if err != nil {
			log.Printf("Error calling Compute on worker at %s: %v", address, err)
			continue
		}
		// Successfully sent the task to a worker
		log.Printf("Worker at %s returned: Result=%v, Rows=%d, Cols=%d, Message=%s", address, res.Result, res.Rows, res.Cols, res.Message)
		return nil
	}

	return errors.New("no available workers to handle the request")
}
