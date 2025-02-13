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
	"sort"
	"strings"
	"sync"
	"time"
)

type Args struct {
	Operation string
	MatrixA   [][]int
	MatrixB   [][]int
}

type Coordinator struct {
	workers     []string
	workerLoads map[string]int // tracks the number of active requests per worker
	mu          sync.Mutex
}

// loadTLSConfig loads the CA certificate and coordinator's certificate/key to configure TLS.
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

// validateMatrixOperation checks that the matrices are dimensionally valid for the requested operation.
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

// callWorkerRPC dials the worker over TLS, makes the RPC call, and returns the result.
func callWorkerRPC(worker string, args *Args) ([][]int, error) {
	// Load the CA certificate (for simplicity, we load it here)
	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	config := &tls.Config{RootCAs: caCertPool}
	conn, err := tls.Dial("tcp", worker, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to worker %s: %v", worker, err)
	}
	defer conn.Close()

	client := rpc.NewClient(conn)
	defer client.Close()

	var result [][]int
	err = client.Call("Worker.PerformOperation", args, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Compute sends the matrix operation to an available worker.
// It tracks worker load and sorts candidates from least loaded to most loaded.
// If all workers fail in the first attempt, the coordinator waits for 1 minute
// before retrying once more. If the second attempt also fails, it returns an error.
func (c *Coordinator) Compute(args *Args, reply *[][]int) error {
	// Validate matrices before sending them to a worker.
	if err := validateMatrixOperation(args); err != nil {
		return err
	}

	var lastErr error

	// We'll attempt at most two full passes.
	for attempt := 1; attempt <= 2; attempt++ {
		// Copy and sort the worker list based on current load.
		c.mu.Lock()
		candidates := make([]string, len(c.workers))
		copy(candidates, c.workers)
		sort.Slice(candidates, func(i, j int) bool {
			return c.workerLoads[candidates[i]] < c.workerLoads[candidates[j]]
		})
		c.mu.Unlock()

		for _, worker := range candidates {
			// Mark the worker as busy by incrementing its load.
			c.mu.Lock()
			c.workerLoads[worker]++
			c.mu.Unlock()

			log.Printf("Forwarding request (%s) to worker: %s (current load: %d)", args.Operation, worker, c.workerLoads[worker])
			result, err := callWorkerRPC(worker, args)

			// Decrement load after the RPC call returns.
			c.mu.Lock()
			c.workerLoads[worker]--
			c.mu.Unlock()

			if err != nil {
				log.Printf("Worker %s error: %v", worker, err)
				lastErr = err
				continue // try next worker
			}

			// If the call succeeded, return the result.
			*reply = result
			return nil
		}

		// All workers failed on this attempt.
		if attempt == 1 {
			var timeSleep = 30 * time.Second
			log.Printf("All workers failed in attempt %d: %v. Waiting for %d seconds before retrying...", attempt, lastErr, int(timeSleep.Seconds()))
			time.Sleep(timeSleep)
		}
	}

	return fmt.Errorf("all workers failed after retrying: %v", lastErr)
}

func main() {
	workers := flag.String("workers", "localhost:1235,localhost:1236,localhost:1237", "Comma-separated list of worker addresses")
	flag.Parse()

	workerList := strings.Split(*workers, ",")

	coordinator := &Coordinator{
		workers:     workerList,
		workerLoads: make(map[string]int),
	}
	// Initialize the load for each worker to zero.
	for _, w := range workerList {
		coordinator.workerLoads[w] = 0
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
