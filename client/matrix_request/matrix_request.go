package matrix_request

// MatrixRequest represents a request containing two matrices and operation type.
type MatrixRequest struct {
	Operation string
	MatrixA   []int
	MatrixB   []int
	RowsA     int
	ColsA     int
	RowsB     int
	ColsB     int
}

// MatrixResponse contains the result of the computation or an error message.
type MatrixResponse struct {
	Message string
}
