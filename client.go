// package main

// import (
// 	"fmt"
// 	"net/rpc"
// )

// func main() {
// 	client, err := rpc.Dial("tcp", "localhost:1234")
// 	if err != nil {
// 		fmt.Println("Error connecting to coordinator:", err)
// 		return
// 	}
// 	defer client.Close()

// 	op := MatrixOperation{
// 		Operation: "Addition",
// 		MatrixA:   [][]int{{1, 2}, {3, 4}},
// 		MatrixB:   [][]int{{5, 6}, {7, 8}},
// 	}

// 	var result Result
// 	err = client.Call("Coordinator.PerformOperation", op, &result)
// 	if err != nil {
// 		fmt.Println("Error performing operation:", err)
// 		return
// 	}

// 	fmt.Println("Result:", result.Matrix)
// }

// package main

// import (
// 	"bufio"
// 	"crypto/tls"
// 	"crypto/x509"
// 	"fmt"
// 	"log"
// 	"net/rpc"
// 	"os"
// 	"strconv"
// 	"strings"
// )

// type Args struct {
// 	Operation string
// 	MatrixA   [][]int
// 	MatrixB   [][]int
// }

// func main() {
// 	// Load CA certificate
// 	caCert, err := os.ReadFile("ca.crt")
// 	if err != nil {
// 		log.Fatal("Failed to read CA certificate:", err)
// 	}

// 	caCertPool := x509.NewCertPool()
// 	caCertPool.AppendCertsFromPEM(caCert)

// 	// Create TLS configuration
// 	config := &tls.Config{
// 		RootCAs: caCertPool,
// 	}

// 	// Dial the coordinator with TLS
// 	conn, err := tls.Dial("tcp", "localhost:1234", config)
// 	if err != nil {
// 		log.Fatal("Failed to dial coordinator:", err)
// 	}
// 	defer conn.Close()

// 	client := rpc.NewClient(conn)
// 	defer client.Close()

// 	reader := bufio.NewReader(os.Stdin)

// 	// Input for Matrix A
// 	fmt.Println("Enter the number of rows for Matrix A:")
// 	rowsA, _ := reader.ReadString('\n')
// 	rowsA = strings.TrimSpace(rowsA)
// 	rowA, _ := strconv.Atoi(rowsA)

// 	fmt.Println("Enter the number of columns for Matrix A:")
// 	colsA, _ := reader.ReadString('\n')
// 	colsA = strings.TrimSpace(colsA)
// 	colA, _ := strconv.Atoi(colsA)

// 	matrixA := make([][]int, rowA)
// 	fmt.Println("Enter the elements of Matrix A (row-wise):")
// 	for i := 0; i < rowA; i++ {
// 		matrixA[i] = make([]int, colA)
// 		for j := 0; j < colA; j++ {
// 			fmt.Printf("Enter element [%d][%d]: ", i, j)
// 			element, _ := reader.ReadString('\n')
// 			element = strings.TrimSpace(element)
// 			matrixA[i][j], _ = strconv.Atoi(element)
// 		}
// 	}

// 	// Input for Matrix B
// 	fmt.Println("Enter the number of rows for Matrix B:")
// 	rowsB, _ := reader.ReadString('\n')
// 	rowsB = strings.TrimSpace(rowsB)
// 	rowB, _ := strconv.Atoi(rowsB)

// 	fmt.Println("Enter the number of columns for Matrix B:")
// 	colsB, _ := reader.ReadString('\n')
// 	colsB = strings.TrimSpace(colsB)
// 	colB, _ := strconv.Atoi(colsB)

// 	matrixB := make([][]int, rowB)
// 	fmt.Println("Enter the elements of Matrix B (row-wise):")
// 	for i := 0; i < rowB; i++ {
// 		matrixB[i] = make([]int, colB)
// 		for j := 0; j < colB; j++ {
// 			fmt.Printf("Enter element [%d][%d]: ", i, j)
// 			element, _ := reader.ReadString('\n')
// 			element = strings.TrimSpace(element)
// 			matrixB[i][j], _ = strconv.Atoi(element)
// 		}
// 	}

// 	// Ask for the operation
// 	fmt.Println("Choose the operation (add, transpose, multiply):")
// 	operation, _ := reader.ReadString('\n')
// 	operation = strings.TrimSpace(operation)

// 	args := Args{
// 		Operation: operation,
// 		MatrixA:   matrixA,
// 		MatrixB:   matrixB,
// 	}

// 	var reply [][]int
// 	err = client.Call("Coordinator.Compute", args, &reply)
// 	if err != nil {
// 		log.Fatal("Coordinator error:", err)
// 	}

// 	fmt.Println("Result:")
// 	for _, row := range reply {
// 		fmt.Println(row)
// 	}
// }

// package main

// import (
// 	"bufio"
// 	"crypto/tls"
// 	"crypto/x509"
// 	"fmt"
// 	"log"
// 	"net/rpc"
// 	"os"
// 	"strconv"
// 	"strings"
// )

// type Args struct {
// 	Operation string
// 	MatrixA   [][]int
// 	MatrixB   [][]int
// }

// func main() {
// 	// Load CA certificate
// 	caCert, err := os.ReadFile("ca.crt")
// 	if err != nil {
// 		log.Fatal("Failed to read CA certificate:", err)
// 	}

// 	caCertPool := x509.NewCertPool()
// 	caCertPool.AppendCertsFromPEM(caCert)

// 	// Create TLS configuration
// 	config := &tls.Config{
// 		RootCAs: caCertPool,
// 		// InsecureSkipVerify: true, // Not recommended for production
// 	}

// 	// Dial the coordinator with TLS
// 	conn, err := tls.Dial("tcp", "localhost:1234", config)
// 	if err != nil {
// 		log.Fatal("Failed to dial coordinator:", err)
// 	}
// 	defer conn.Close()

// 	client := rpc.NewClient(conn)
// 	defer client.Close()

// 	reader := bufio.NewReader(os.Stdin)

// 	// Input for Matrix A
// 	fmt.Println("Enter the number of rows for Matrix A:")
// 	rowsA, _ := reader.ReadString('\n')
// 	rowsA = strings.TrimSpace(rowsA)
// 	rowA, _ := strconv.Atoi(rowsA)

// 	fmt.Println("Enter the number of columns for Matrix A:")
// 	colsA, _ := reader.ReadString('\n')
// 	colsA = strings.TrimSpace(colsA)
// 	colA, _ := strconv.Atoi(colsA)

// 	matrixA := make([][]int, rowA)
// 	fmt.Println("Enter the elements of Matrix A (row-wise):")
// 	for i := 0; i < rowA; i++ {
// 		matrixA[i] = make([]int, colA)
// 		for j := 0; j < colA; j++ {
// 			fmt.Printf("Enter element [%d][%d]: ", i, j)
// 			element, _ := reader.ReadString('\n')
// 			element = strings.TrimSpace(element)
// 			matrixA[i][j], _ = strconv.Atoi(element)
// 		}
// 	}

// 	// Input for Matrix B
// 	fmt.Println("Enter the number of rows for Matrix B:")
// 	rowsB, _ := reader.ReadString('\n')
// 	rowsB = strings.TrimSpace(rowsB)
// 	rowB, _ := strconv.Atoi(rowsB)

// 	fmt.Println("Enter the number of columns for Matrix B:")
// 	colsB, _ := reader.ReadString('\n')
// 	colsB = strings.TrimSpace(colsB)
// 	colB, _ := strconv.Atoi(colsB)

// 	matrixB := make([][]int, rowB)
// 	fmt.Println("Enter the elements of Matrix B (row-wise):")
// 	for i := 0; i < rowB; i++ {
// 		matrixB[i] = make([]int, colB)
// 		for j := 0; j < colB; j++ {
// 			fmt.Printf("Enter element [%d][%d]: ", i, j)
// 			element, _ := reader.ReadString('\n')
// 			element = strings.TrimSpace(element)
// 			matrixB[i][j], _ = strconv.Atoi(element)
// 		}
// 	}

// 	// Ask for the operation
// 	fmt.Println("Choose the operation (add, transpose, multiply):")
// 	operation, _ := reader.ReadString('\n')
// 	operation = strings.TrimSpace(operation)

// 	args := Args{
// 		Operation: operation,
// 		MatrixA:   matrixA,
// 		MatrixB:   matrixB,
// 	}

// 	var reply [][]int
// 	err = client.Call("Coordinator.Compute", args, &reply)
// 	if err != nil {
// 		log.Fatal("Coordinator error:", err)
// 	}

// 	fmt.Println("Result:")
// 	for _, row := range reply {
// 		fmt.Println(row)
// 	}
// }

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
