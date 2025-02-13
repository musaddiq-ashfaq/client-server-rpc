package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
	strings "strings"

	mr "client/matrix_request"
)

// printMatrix prints a 1D slice as a 2D matrix.
func printMatrix(matrix []int, rows, cols int) {
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			fmt.Printf("%d ", matrix[i*cols+j])
		}
		fmt.Println()
	}
}

// getMatrixInput takes user input for a matrix and validates it.
func getMatrixInput(rows, cols int) []int {
	matrix := make([]int, rows*cols)
	fmt.Printf("Enter %d elements (space-separated) for a %dx%d matrix:\n", rows*cols, rows, cols)

	reader := bufio.NewReader(os.Stdin)
	for {
		input, _ := reader.ReadString('\n')
		nums := strings.Fields(input)

		if len(nums) != rows*cols {
			fmt.Printf("Invalid input. Please enter exactly %d numbers:\n", rows*cols)
			continue
		}

		valid := true
		for i, num := range nums {
			val, err := strconv.Atoi(num)
			if err != nil {
				fmt.Println("Invalid number, please enter integers only.")
				valid = false
				break
			}
			matrix[i] = val
		}
		if valid {
			break
		}
	}
	return matrix
}

func main() {
	// Establish a TLS connection
	config := &tls.Config{InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", "localhost:50051", config)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := rpc.NewClient(conn)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Choose operation (add, multiply, transpose):")
	operation, _ := reader.ReadString('\n')
	operation = strings.TrimSpace(operation)

	var req mr.MatrixRequest
	req.Operation = operation

	fmt.Println("Enter number of rows and columns for Matrix A:")
	var rowsA, colsA int
	fmt.Scan(&rowsA, &colsA)
	req.MatrixA = getMatrixInput(rowsA, colsA)
	req.RowsA = rowsA
	req.ColsA = colsA

	if operation == "add" || operation == "multiply" {
		fmt.Println("Enter number of rows and columns for Matrix B:")
		var rowsB, colsB int
		fmt.Scan(&rowsB, &colsB)

		if operation == "add" && (rowsA != rowsB || colsA != colsB) {
			log.Fatal("For addition, both matrices must have the same dimensions.")
		}
		if operation == "multiply" && colsA != rowsB {
			log.Fatal("For multiplication, columns of A must match rows of B.")
		}

		req.MatrixB = getMatrixInput(rowsB, colsB)
		req.RowsB = rowsB
		req.ColsB = colsB
	}

	var res mr.MatrixResponse
	err = client.Call("MatrixService.Serve", req, &res)
	if err != nil {
		log.Fatalf("RPC error: %v", err)
	}

	fmt.Println("Response from coordinator:", res.Message)
	fmt.Println("Result Matrix:")
	printMatrix(res.Result, res.Rows, res.Cols)
}
