// package main

// import (
// 	"fmt"
// 	"log"
// 	"net"
// 	"net/rpc"
// 	"sync"
// )

// type MatrixOperation struct {
// 	Operation string
// 	MatrixA   [][]int
// 	MatrixB   [][]int
// }

// type Result struct {
// 	Matrix [][]int
// 	Error  string
// }

// type Worker struct {
// 	ID     int
// 	Busy   bool
// 	Client *rpc.Client
// }

// var (
// 	workers []*Worker
// 	mutex   sync.Mutex
// )

// type Coordinator struct{}

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

// 	err := selectedWorker.Client.Call("Worker.Compute", op, result)

// 	mutex.Lock()
// 	selectedWorker.Busy = false
// 	mutex.Unlock()

// 	return err
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

// package main

// import (
// 	"crypto/tls"
// 	"crypto/x509"
// 	"errors"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"net/rpc"
// 	"os"
// 	"sync"
// )

// type Args struct {
// 	Operation string
// 	MatrixA   [][]int
// 	MatrixB   [][]int
// }

// type Coordinator struct {
// 	workers []string
// 	mu      sync.Mutex
// }

// func (c *Coordinator) Compute(args *Args, reply *[][]int) error {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()

// 	for _, worker := range c.workers {
// 		// Load CA certificate
// 		caCert, err := os.ReadFile("ca.crt")
// 		if err != nil {
// 			log.Fatal("Failed to read CA certificate:", err)
// 		}

// 		caCertPool := x509.NewCertPool()
// 		caCertPool.AppendCertsFromPEM(caCert)

// 		// Create TLS configuration
// 		config := &tls.Config{
// 			RootCAs: caCertPool,
// 		}

// 		// Dial the worker with TLS
// 		conn, err := tls.Dial("tcp", worker, config)
// 		if err != nil {
// 			log.Printf("Worker %s failed: %v", worker, err)
// 			continue
// 		}
// 		defer conn.Close()

// 		client := rpc.NewClient(conn)
// 		defer client.Close()

// 		var result [][]int
// 		err = client.Call("Worker.PerformOperation", args, &result)
// 		if err != nil {
// 			log.Printf("Worker %s failed to perform operation: %v", worker, err)
// 			continue
// 		}

// 		*reply = result
// 		return nil
// 	}

// 	return errors.New("all workers failed")
// }

// func main() {
// 	coordinator := &Coordinator{
// 		workers: []string{"localhost:1235", "localhost:1236", "localhost:1237"},
// 	}

// 	rpc.Register(coordinator)
// 	rpc.HandleHTTP()

// 	// Load server certificate and key
// 	cert, err := tls.LoadX509KeyPair("localhost.crt", "localhost.key")
// 	if err != nil {
// 		log.Fatal("Failed to load server certificate:", err)
// 	}

// 	// Create TLS configuration
// 	config := &tls.Config{
// 		Certificates: []tls.Certificate{cert},
// 	}

// 	// Listen for incoming connections with TLS
// 	listener, err := tls.Listen("tcp", ":1234", config)
// 	if err != nil {
// 		log.Fatal("listen error:", err)
// 	}

// 	fmt.Println("Coordinator is running...")
// 	http.Serve(listener, nil)
// }

// package main

// import (
// 	"crypto/tls"
// 	"crypto/x509"
// 	"errors"
// 	"flag"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"net/rpc"
// 	"os"
// 	"strings"
// 	"sync"
// )

// type Args struct {
// 	Operation string
// 	MatrixA   [][]int
// 	MatrixB   [][]int
// }

// type Coordinator struct {
// 	workers []string
// 	mu      sync.Mutex
// }

// func (c *Coordinator) Compute(args *Args, reply *[][]int) error {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()

// 	for _, worker := range c.workers {
// 		// Load CA certificate
// 		caCert, err := os.ReadFile("ca.crt")
// 		if err != nil {
// 			log.Fatal("Failed to read CA certificate:", err)
// 		}

// 		caCertPool := x509.NewCertPool()
// 		caCertPool.AppendCertsFromPEM(caCert)

// 		// Create TLS configuration
// 		config := &tls.Config{
// 			RootCAs: caCertPool,
// 		}

// 		// Dial the worker with TLS
// 		conn, err := tls.Dial("tcp", worker, config)
// 		if err != nil {
// 			log.Printf("Worker %s failed: %v", worker, err)
// 			continue
// 		}
// 		defer conn.Close()

// 		client := rpc.NewClient(conn)
// 		defer client.Close()

// 		var result [][]int
// 		err = client.Call("Worker.PerformOperation", args, &result)
// 		if err != nil {
// 			log.Printf("Worker %s failed to perform operation: %v", worker, err)
// 			continue
// 		}

// 		*reply = result
// 		return nil
// 	}

// 	return errors.New("all workers failed")
// }

// func main() {
// 	// Define command-line flags for worker addresses
// 	workers := flag.String("workers", "localhost:1235,localhost:1236,localhost:1237", "Comma-separated list of worker addresses")
// 	flag.Parse()

// 	coordinator := &Coordinator{
// 		workers: strings.Split(*workers, ","),
// 	}

// 	rpc.Register(coordinator)
// 	rpc.HandleHTTP()

// 	// Load server certificate and key
// 	cert, err := tls.LoadX509KeyPair("coordinator.crt", "coordinator.key")
// 	if err != nil {
// 		log.Fatal("Failed to load server certificate:", err)
// 	}

// 	// Create TLS configuration for one-way TLS
// 	config := &tls.Config{
// 		Certificates: []tls.Certificate{cert}, // Coordinator's certificate
// 		ClientAuth:   tls.NoClientCert,        // Do not require a client certificate
// 	}
// 	// Listen for incoming connections with TLS
// 	listener, err := tls.Listen("tcp", ":1234", config)
// 	if err != nil {
// 		log.Fatal("listen error:", err)
// 	}

// 	fmt.Println("Coordinator is running...")
// 	http.Serve(listener, nil)
// }

// package main

// import (
// 	"crypto/tls"
// 	"crypto/x509"
// 	"errors"
// 	"flag"
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"net/rpc"
// 	"os"
// 	"strings"
// 	"sync"
// 	"time"
// )

// type Args struct {
// 	Operation string
// 	MatrixA   [][]int
// 	MatrixB   [][]int
// }

// type Response struct {
// 	Result [][]int
// 	Error  string // To capture errors like dimension mismatches
// }

// type Coordinator struct {
// 	workers []string
// 	mu      sync.Mutex
// }

// // Selects a random available worker
// func (c *Coordinator) selectWorker() (string, error) {
// 	if len(c.workers) == 0 {
// 		return "", errors.New("no available workers")
// 	}
// 	rand.Seed(time.Now().UnixNano()) // Seed randomizer
// 	return c.workers[rand.Intn(len(c.workers))], nil
// }

// // Establish TLS configuration
// func loadTLSConfig() (*tls.Config, error) {
// 	// Load CA certificate
// 	caCert, err := os.ReadFile("ca.crt")
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read CA certificate: %v", err)
// 	}

// 	caCertPool := x509.NewCertPool()
// 	caCertPool.AppendCertsFromPEM(caCert)

// 	// Load server certificate and key
// 	cert, err := tls.LoadX509KeyPair("coordinator.crt", "coordinator.key")
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to load server certificate: %v", err)
// 	}

// 	return &tls.Config{
// 		Certificates: []tls.Certificate{cert}, // Coordinator's certificate
// 		RootCAs:      caCertPool,
// 		ClientAuth:   tls.NoClientCert, // No client certificate required
// 	}, nil
// }

// // Compute matrix operation using a worker
// func (c *Coordinator) Compute(args *Args, reply *Response) error {
// 	c.mu.Lock()
// 	defer c.mu.Unlock()

// 	fmt.Println("Received request: Operation =", args.Operation) // ✅ Debug log

// 	for _, worker := range c.workers {
// 		fmt.Println("Trying worker:", worker) // ✅ Debug log

// 		// Load CA certificate
// 		caCert, err := os.ReadFile("ca.crt")
// 		if err != nil {
// 			log.Fatal("Failed to read CA certificate:", err)
// 		}

// 		caCertPool := x509.NewCertPool()
// 		caCertPool.AppendCertsFromPEM(caCert)

// 		// Create TLS configuration
// 		config := &tls.Config{
// 			RootCAs: caCertPool,
// 		}

// 		// Dial the worker with TLS
// 		conn, err := tls.Dial("tcp", worker, config)
// 		if err != nil {
// 			log.Printf("Failed to connect to worker %s: %v", worker, err)
// 			continue
// 		}
// 		defer conn.Close()

// 		client := rpc.NewClient(conn)
// 		defer client.Close()

// 		fmt.Println("Connected to worker:", worker) // ✅ Debug log

// 		var result Response
// 		err = client.Call("Worker.PerformOperation", args, &result)
// 		if err != nil {
// 			log.Printf("Worker %s failed to perform operation: %v", worker, err)
// 			continue
// 		}

// 		// Check if worker returned an error (like a matrix size mismatch)
// 		if result.Error != "" {
// 			fmt.Println("Worker error:", result.Error)
// 			*reply = result // Send the error to the client
// 			return nil
// 		}

// 		fmt.Println("Worker", worker, "successfully computed result") // ✅ Debug log

// 		*reply = result
// 		return nil
// 	}

// 	return errors.New("all workers failed")
// }

// func main() {
// 	// Define command-line flags for worker addresses
// 	workers := flag.String("workers", "localhost:1235,localhost:1236,localhost:1237", "Comma-separated list of worker addresses")
// 	flag.Parse()

// 	coordinator := &Coordinator{
// 		workers: strings.Split(*workers, ","),
// 	}

// 	err := rpc.Register(coordinator)
// 	if err != nil {
// 		log.Fatalf("Error registering RPC service: %v", err)
// 	}

// 	config, err := loadTLSConfig()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	listener, err := tls.Listen("tcp", ":1234", config)
// 	if err != nil {
// 		log.Fatal("Listen error:", err)
// 	}

// 	fmt.Println("Coordinator is running...")
// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			log.Println("Connection error:", err)
// 			continue
// 		}
// 		go rpc.ServeConn(conn)
// 	}
// }

package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"
	"sync"
)

type Args struct {
	Operation string
	MatrixA   [][]int
	MatrixB   [][]int
}

type Coordinator struct {
	workers []string
	mu      sync.Mutex
}

// Establish TLS configuration for incoming client connections.
func loadTLSConfig() (*tls.Config, error) {
	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair("coordinator.crt", "coordinator.key")
	if err != nil {
		return nil, fmt.Errorf("failed to load server certificate: %v", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		ClientAuth:   tls.NoClientCert,
	}, nil
}

// Validate matrix dimensions before sending them to the worker.
func validateMatrixOperation(args *Args) error {
	switch args.Operation {
	case "add":
		if len(args.MatrixA) != len(args.MatrixB) || len(args.MatrixA[0]) != len(args.MatrixB[0]) {
			return errors.New("matrix addition error: matrices must have the same dimensions")
		}
	case "multiply":
		if len(args.MatrixA[0]) != len(args.MatrixB) {
			return errors.New("matrix multiplication error: columns of A must match rows of B")
		}
	case "transpose":
		if len(args.MatrixB) > 0 {
			return errors.New("transpose operation should only have one matrix")
		}
	default:
		return errors.New("invalid operation")
	}
	return nil
}

// Compute sends the matrix operation to an available worker.
// If a connection error occurs or the worker returns an error, it will try the next worker.
// If all workers fail, it returns an error.
func (c *Coordinator) Compute(args *Args, reply *[][]int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	fmt.Println("Received request: Operation =", args.Operation)

	// Validate matrices before sending them to a worker.
	if err := validateMatrixOperation(args); err != nil {
		return err
	}

	// Load the CA certificate once outside the loop.
	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
		return fmt.Errorf("failed to read CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	var lastErr error

	for _, worker := range c.workers {
		fmt.Println("Trying worker:", worker)

		config := &tls.Config{RootCAs: caCertPool}

		conn, err := tls.Dial("tcp", worker, config)
		if err != nil {
			log.Printf("Failed to connect to worker %s: %v", worker, err)
			lastErr = err
			continue
		}

		client := rpc.NewClient(conn)
		var result [][]int
		err = client.Call("Worker.PerformOperation", args, &result)
		if err != nil {
			log.Printf("Worker %s failed to perform operation: %v", worker, err)
			lastErr = err
			client.Close()
			conn.Close()
			continue
		}

		// Close the client and connection before returning.
		client.Close()
		conn.Close()
		*reply = result
		return nil
	}

	return fmt.Errorf("all workers failed: %v", lastErr)
}

func main() {
	workers := flag.String("workers", "localhost:1235,localhost:1236,localhost:1237", "Comma-separated list of worker addresses")
	flag.Parse()

	coordinator := &Coordinator{
		workers: strings.Split(*workers, ","),
	}

	err := rpc.Register(coordinator)
	if err != nil {
		log.Fatalf("Error registering RPC service: %v", err)
	}

	config, err := loadTLSConfig()
	if err != nil {
		log.Fatal(err)
	}

	listener, err := tls.Listen("tcp", ":1234", config)
	if err != nil {
		log.Fatal("Listen error:", err)
	}

	fmt.Println("Coordinator is running...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
