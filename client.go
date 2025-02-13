package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

type Args struct {
	Operation string
	MatrixA   [][]int
	MatrixB   [][]int
}

func main() {
	// Accept coordinator address from command line
	coordinatorAddr := flag.String("coordinator", "localhost:1234", "Coordinator address (IP:port)")
	flag.Parse()

	// Load CA certificate
	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
		log.Fatal("Failed to read CA certificate:", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create TLS configuration
	config := &tls.Config{
		RootCAs: caCertPool,
	}

	// Connect to the coordinator
	conn, err := tls.Dial("tcp", *coordinatorAddr, config)
	if err != nil {
		log.Fatal("Failed to connect to coordinator:", err)
	}
	defer conn.Close()

	client := rpc.NewClient(conn)
	defer client.Close()

	reader := bufio.NewReader(os.Stdin)

	// Input for Matrix A
	fmt.Println("Enter the number of rows for Matrix A:")
	rowA, _ := strconv.Atoi(strings.TrimSpace(readInput(reader)))

	fmt.Println("Enter the number of columns for Matrix A:")
	colA, _ := strconv.Atoi(strings.TrimSpace(readInput(reader)))

	matrixA := readMatrix(reader, rowA, colA, "A")

	// Input for Matrix B
	fmt.Println("Enter the number of rows for Matrix B:")
	rowB, _ := strconv.Atoi(strings.TrimSpace(readInput(reader)))

	fmt.Println("Enter the number of columns for Matrix B:")
	colB, _ := strconv.Atoi(strings.TrimSpace(readInput(reader)))

	matrixB := readMatrix(reader, rowB, colB, "B")

	// Ask for the operation
	fmt.Println("Choose the operation (add, transpose, multiply):")
	operation := strings.TrimSpace(readInput(reader))

	args := Args{
		Operation: operation,
		MatrixA:   matrixA,
		MatrixB:   matrixB,
	}

	var reply [][]int
	err = client.Call("Coordinator.Compute", args, &reply)
	if err != nil {
		log.Fatal("Coordinator error:", err)
	}

	// Display result
	fmt.Println("Result:")
	for _, row := range reply {
		fmt.Println(row)
	}
}

// Helper function to read input
func readInput(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// Helper function to read matrix input
func readMatrix(reader *bufio.Reader, rows, cols int, name string) [][]int {
	matrix := make([][]int, rows)
	fmt.Printf("Enter the elements of Matrix %s (row-wise):\n", name)
	for i := 0; i < rows; i++ {
		matrix[i] = make([]int, cols)
		for j := 0; j < cols; j++ {
			fmt.Printf("Enter element [%d][%d]: ", i, j)
			matrix[i][j], _ = strconv.Atoi(strings.TrimSpace(readInput(reader)))
		}
	}
	return matrix
}
