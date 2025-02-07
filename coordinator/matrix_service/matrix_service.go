package matrix_service

import (
	"errors"
	"log"
	client "client/matrix_request"
)

// MatrixService defines an RPC service with a Compute method.
type MatrixService struct{}

// Compute handles matrix operations (for now, it just acknowledges the request).
func (m *MatrixService) Compute(req client.MatrixRequest, res *client.MatrixResponse) error {
	if len(req.MatrixA) == 0 || len(req.MatrixB) == 0 {
		return errors.New("invalid matrices")
	}
	log.Printf("Received request: Operation=%v, MatrixA=%v, MatrixB=%v\n", req.Operation, req.MatrixA, req.MatrixB)
	res.Message = "Request received by coordinator"
	return nil
}
