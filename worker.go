// package main

// import (
// 	"fmt"
// 	"net"
// 	"net/rpc"
// )

// type WorkerService struct{}

// func (w *WorkerService) Compute(op MatrixOperation, result *Result) error {
// 	var res [][]int
// 	switch op.Operation {
// 	case "Addition":
// 		res = addMatrices(op.MatrixA, op.MatrixB)
// 	case "Transpose":
// 		res = transposeMatrix(op.MatrixA)
// 	case "Multiplication":
// 		res = multiplyMatrices(op.MatrixA, op.MatrixB)
// 	default:
// 		result.Error = "Invalid operation"
// 		return nil
// 	}
// 	result.Matrix = res
// 	return nil
// }

// func main() {
// 	worker := new(WorkerService)
// 	rpc.Register(worker)
// 	listener, err := net.Listen("tcp", ":1235")
// 	if err != nil {
// 		fmt.Println("Error starting worker:", err)
// 		return
// 	}
// 	defer listener.Close()

// 	fmt.Println("Worker server is running on port 1235...")
// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			fmt.Println("Connection error:", err)
// 			continue
// 		}
// 		go rpc.ServeConn(conn)
// 	}
// }

// package main

// import (
// 	"crypto/tls"
// 	"errors"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"net/rpc"
// )

// type Args struct {
// 	Operation string
// 	MatrixA   [][]int
// 	MatrixB   [][]int
// }

// type Worker struct{}

// func (w *Worker) PerformOperation(args *Args, reply *[][]int) error {
// 	switch args.Operation {
// 	case "add":
// 		*reply = addMatrices(args.MatrixA, args.MatrixB)
// 	case "transpose":
// 		*reply = transposeMatrix(args.MatrixA)
// 	case "multiply":
// 		*reply = multiplyMatrices(args.MatrixA, args.MatrixB)
// 	default:
// 		return errors.New("unsupported operation")
// 	}
// 	return nil
// }

// func addMatrices(a, b [][]int) [][]int {
// 	result := make([][]int, len(a))
// 	for i := range a {
// 		result[i] = make([]int, len(a[i]))
// 		for j := range a[i] {
// 			result[i][j] = a[i][j] + b[i][j]
// 		}
// 	}
// 	return result
// }

// func transposeMatrix(a [][]int) [][]int {
// 	result := make([][]int, len(a[0]))
// 	for i := range result {
// 		result[i] = make([]int, len(a))
// 		for j := range a {
// 			result[i][j] = a[j][i]
// 		}
// 	}
// 	return result
// }

// func multiplyMatrices(a, b [][]int) [][]int {
// 	result := make([][]int, len(a))
// 	for i := range a {
// 		result[i] = make([]int, len(b[0]))
// 		for j := range b[0] {
// 			for k := range b {
// 				result[i][j] += a[i][k] * b[k][j]
// 			}
// 		}
// 	}
// 	return result
// }

// func main() {
// 	worker := new(Worker)
// 	rpc.Register(worker)
// 	rpc.HandleHTTP()

// 	// Load server certificate and key
// 	cert, err := tls.LoadX509KeyPair("localhost.crt", "localhost.key") // Change for each worker
// 	if err != nil {
// 		log.Fatal("Failed to load server certificate:", err)
// 	}

// 	// Create TLS configuration
// 	config := &tls.Config{
// 		Certificates: []tls.Certificate{cert},
// 	}

// 	// Listen for incoming connections with TLS
// 	listener, err := tls.Listen("tcp", ":1235", config) // Change port for each worker
// 	if err != nil {
// 		log.Fatal("listen error:", err)
// 	}

// 	fmt.Println("Worker is running...")
// 	http.Serve(listener, nil)
// }

package main

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/rpc"
	"os"
)

type Args struct {
	Operation string
	MatrixA   [][]int
	MatrixB   [][]int
}

type Worker struct{}

// PerformOperation executes the requested matrix operation and returns the result.
func (w *Worker) PerformOperation(args *Args, reply *[][]int) error {
	log.Println("Received operation:", args.Operation)

	// Check the operation type
	switch args.Operation {
	case "add":
		*reply = addMatrices(args.MatrixA, args.MatrixB)
	case "multiply":
		*reply = multiplyMatrices(args.MatrixA, args.MatrixB)
	case "transpose":
		*reply = transposeMatrix(args.MatrixA)
	default:
		return errors.New("invalid operation")
	}

	log.Println("Computed result:", *reply)
	return nil
}

// Matrix Addition
func addMatrices(a, b [][]int) [][]int {
	// Ensure both matrices have the same dimensions
	if len(a) != len(b) || len(a[0]) != len(b[0]) {
		log.Fatal("Matrix dimensions do not match for addition")
	}

	result := make([][]int, len(a))
	for i := range a {
		if len(a[i]) != len(b[i]) { // Check row length consistency
			log.Fatal("Row size mismatch in matrices")
		}
		result[i] = make([]int, len(a[i]))
		for j := range a[i] {
			result[i][j] = a[i][j] + b[i][j]
		}
	}
	return result
}

// Matrix Transposition
func transposeMatrix(a [][]int) [][]int {
	// Ensure matrix is non-empty
	if len(a) == 0 || len(a[0]) == 0 {
		log.Fatal("Matrix is empty, cannot transpose")
	}

	result := make([][]int, len(a[0]))
	for i := range result {
		result[i] = make([]int, len(a))
		for j := range a {
			if i >= len(a[j]) { // Ensure valid index access
				log.Fatal("Matrix is not rectangular, cannot transpose")
			}
			result[i][j] = a[j][i]
		}
	}
	return result
}

// Matrix Multiplication
func multiplyMatrices(a, b [][]int) [][]int {
	// Ensure matrix dimensions allow multiplication
	if len(a) == 0 || len(b) == 0 || len(a[0]) != len(b) {
		log.Fatal("Matrix dimensions do not allow multiplication")
	}

	result := make([][]int, len(a))
	for i := range a {
		result[i] = make([]int, len(b[0]))
		for j := range b[0] {
			for k := range b {
				if k >= len(a[i]) || j >= len(b[k]) { // Check bounds
					log.Fatal("Matrix dimension mismatch during multiplication")
				}
				result[i][j] += a[i][k] * b[k][j]
			}
		}
	}
	return result
}

func main() {
	// Define command-line flags
	port := flag.String("port", "1235", "Port for the worker to listen on")
	certFile := flag.String("cert", "worker1.crt", "Path to the certificate file")
	keyFile := flag.String("key", "worker1.key", "Path to the private key file")
	flag.Parse()

	// Check if certificate and key files exist
	if _, err := os.Stat(*certFile); os.IsNotExist(err) {
		log.Fatalf("Certificate file %s does not exist", *certFile)
	}
	if _, err := os.Stat(*keyFile); os.IsNotExist(err) {
		log.Fatalf("Private key file %s does not exist", *keyFile)
	}

	// Register the RPC worker
	worker := new(Worker)
	err := rpc.Register(worker)
	if err != nil {
		log.Fatal("Error registering worker:", err)
	}

	// Load server certificate and key
	cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
	if err != nil {
		log.Fatal("Failed to load server certificate:", err)
	}

	// Create TLS configuration
	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	// Listen for incoming TLS connections
	listener, err := tls.Listen("tcp", ":"+*port, config)
	if err != nil {
		log.Fatal("Listen error:", err)
	}

	fmt.Printf("Worker is running on port %s...\n", *port)
	log.Println("Worker registered and ready to receive RPC calls...")

	// Accept connections and serve RPC requests
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn) // Handle each connection in a new goroutine
	}
}
