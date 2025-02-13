# Blockchain Cryptocurrency Client-Server RPC-based Matrix Computation

This project implements a **Blockchain-based Cryptocurrency** system with a **Client-Server architecture** using **Remote Procedure Calls (RPC)** for **secure and decentralized matrix computation**. The system ensures **secure communication** via TLS encryption and involves **a coordinator, client, and worker nodes**.

## Features
- **Blockchain-based cryptocurrency system** with decentralized architecture.
- **Client-Server Model**: The client interacts with a central coordinator and worker nodes.
- **Remote Procedure Calls (RPC)** for distributed matrix computations.
- **Secure TLS communication** with certificate-based authentication.
- **Efficient multi-node computation** using Go.

## Project Structure
```
├── client.go              # Client-side logic
├── coordinator.go         # Central coordinator for task distribution
├── worker.go              # Worker nodes performing computations
├── ca.crt, ca.key         # Certificate Authority (CA) files for TLS security
├── coordinator.crt, key   # TLS certificates for coordinator
├── worker.crt, key        # TLS certificates for worker
├── openssl.cnf            # OpenSSL configuration for certificate generation
├── .gitignore             # Git ignore rules
└── README.md              # Project documentation
```

## Installation & Setup
### Prerequisites
- **Go (Golang)** installed (>= 1.18)
- **OpenSSL** for certificate generation

### Steps
1. **Clone the Repository**
   ```sh
   git clone https://github.com/your-repo/blockchain-rpc.git
   cd blockchain-rpc
   ```

2. **Generate TLS Certificates**
   Run the provided OpenSSL commands to generate necessary certificates:
   ```sh
   openssl req -x509 -newkey rsa:4096 -keyout ca.key -out ca.crt -days 365
   ```
   *(Refer to `Certificate Generation Keys.txt` for detailed steps.)*

3. **Build the Project**
   ```sh
   go build client.go
   go build coordinator.go
   go build worker.go
   ```

4. **Run the Coordinator**
   ```sh
   ./coordinator
   ```

5. **Run the Worker Nodes**
   ```sh
   ./worker
   ```

6. **Run the Client**
   ```sh
   ./client
   ```

## Usage
- The **coordinator** assigns matrix computation tasks to worker nodes.
- The **worker nodes** perform the computations and return results securely.
- The **client** interacts with the coordinator to submit tasks and receive computed results.

## Contributing
1. **Fork the repository**
2. **Create a feature branch**
3. **Commit your changes**
4. **Push to your fork & create a Pull Request**

## License
This project is licensed under the **MIT License**.

## Contact
For queries or contributions, reach out via **[abr78384@gmail.com](mailto:abr78384@gmail.com)**.