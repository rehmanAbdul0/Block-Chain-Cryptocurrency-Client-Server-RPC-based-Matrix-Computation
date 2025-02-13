// // Go implementation of the Blockchain and Cryptocurrency Assignment #1

// package main

// import (
// 	"fmt"
// 	"log"
// 	"net"
// 	"net/rpc"
// 	"sync"
// )

// // MatrixOperation represents a computation request
// type MatrixOperation struct {
// 	Operation string
// 	MatrixA   [][]int
// 	MatrixB   [][]int
// }

// // Result represents the result returned from workers
// type Result struct {
// 	Matrix [][]int
// 	Error  string
// }

// // Worker struct to manage computation
// type Worker struct {
// 	ID   int
// 	Busy bool
// 	Client *rpc.Client
// }

// var (
// 	workers []*Worker
// 	mutex   sync.Mutex
// )

// // Coordinator struct
// type Coordinator struct {}

// // RegisterWorker allows a worker to register itself
// func (c *Coordinator) RegisterWorker(address string, reply *string) error {
// 	mutex.Lock()
// 	defer mutex.Unlock()

// 	client, err := rpc.Dial("tcp", address)
// 	if err != nil {
// 		return err
// 	}

// 	worker := &Worker{ID: len(workers) + 1, Busy: false, Client: client}
// 	workers = append(workers, worker)
// 	*reply = fmt.Sprintf("Worker %d registered successfully", worker.ID)
// 	return nil
// }

// // PerformOperation handles matrix operation requests
// func (c *Coordinator) PerformOperation(op MatrixOperation, result *Result) error {
// 	mutex.Lock()
// 	var selectedWorker *Worker
// 	for _, worker := range workers {
// 		if !worker.Busy {
// 			selectedWorker = worker
// 			worker.Busy = true
// 			break
// 		}
// 	}
// 	mutex.Unlock()

// 	if selectedWorker == nil {
// 		return fmt.Errorf("No available workers")
// 	}

// 	// Perform operation on worker
// 	err := selectedWorker.Client.Call("Worker.Compute", op, result)

// 	mutex.Lock()
// 	selectedWorker.Busy = false
// 	mutex.Unlock()

// 	return err
// }

// // Worker implementation
// type WorkerService struct {}

// // Compute handles matrix operations
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

// // addMatrices performs matrix addition
// func addMatrices(A, B [][]int) [][]int {
// 	rows, cols := len(A), len(A[0])
// 	result := make([][]int, rows)
// 	for i := range result {
// 		result[i] = make([]int, cols)
// 		for j := range result[i] {
// 			result[i][j] = A[i][j] + B[i][j]
// 		}
// 	}
// 	return result
// }

// // transposeMatrix performs matrix transposition
// func transposeMatrix(A [][]int) [][]int {
// 	rows, cols := len(A), len(A[0])
// 	result := make([][]int, cols)
// 	for i := range result {
// 		result[i] = make([]int, rows)
// 		for j := range result[i] {
// 			result[i][j] = A[j][i]
// 		}
// 	}
// 	return result
// }

// // multiplyMatrices performs matrix multiplication
// func multiplyMatrices(A, B [][]int) [][]int {
// 	rowsA, colsA, colsB := len(A), len(A[0]), len(B[0])
// 	result := make([][]int, rowsA)
// 	for i := range result {
// 		result[i] = make([]int, colsB)
// 		for j := range result[i] {
// 			for k := 0; k < colsA; k++ {
// 				result[i][j] += A[i][k] * B[k][j]
// 			}
// 		}
// 	}
// 	return result
// }

// func main() {
// 	coordinator := new(Coordinator)
// 	rpc.Register(coordinator)
// 	listener, err := net.Listen("tcp", ":1234")
// 	if err != nil {
// 		log.Fatal("Error starting server:", err)
// 	}
// 	defer listener.Close()

// 	fmt.Println("Coordinator server is running on port 1234...")
// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			log.Println("Connection error:", err)
// 			continue
// 		}
// 		go rpc.ServeConn(conn)
// 	}
// }