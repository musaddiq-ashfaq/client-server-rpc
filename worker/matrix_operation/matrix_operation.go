package matrix_operation

// Import the shared request/response structures from the client package
import (
	client "client/matrix_request"
	"errors"
)

// PerformMatrixOperation handles matrix operations
func PerformMatrixOperation(req client.MatrixRequest) (client.MatrixResponse, error) {
	var result []int
	var rows, cols int

	switch req.Operation {
	case "add":
		result = matrixAddition(req.MatrixA, req.MatrixB)
		rows, cols = req.RowsA, req.ColsA // Same dimensions as inputs
	case "transpose":
		result = matrixTranspose(req.MatrixA, req.RowsA, req.ColsA)
		rows, cols = req.ColsA, req.RowsA // Transposed dimensions
	case "multiply":
		result = matrixMultiplication(req.MatrixA, req.MatrixB, req.RowsA, req.ColsA, req.ColsB)
		rows, cols = req.RowsA, req.ColsB // Rows from A, Cols from B
	default:
		return client.MatrixResponse{Message: "Invalid operation"}, errors.New("unsupported operation")
	}

	return client.MatrixResponse{
		Result:  result,
		Rows:    rows,
		Cols:    cols,
		Message: "Computation successful",
	}, nil
}

// matrixAddition performs element-wise addition of two matrices.
func matrixAddition(A, B []int) []int {
	result := make([]int, len(A))
	for i := range A {
		result[i] = A[i] + B[i]
	}
	return result
}

// matrixTranspose computes the transpose of a matrix.
func matrixTranspose(A []int, rows, cols int) []int {
	result := make([]int, len(A))
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			result[c*rows+r] = A[r*cols+c]
		}
	}
	return result
}

// matrixMultiplication multiplies two matrices.
func matrixMultiplication(A, B []int, rowsA, colsA, colsB int) []int {
	result := make([]int, rowsA*colsB)
	for i := 0; i < rowsA; i++ {
		for j := 0; j < colsB; j++ {
			sum := 0
			for k := 0; k < colsA; k++ {
				sum += A[i*colsA+k] * B[k*colsB+j]
			}
			result[i*colsB+j] = sum
		}
	}
	return result
}
