# Client-Server-RPC

## Overview
This project implements a client-server system where a client requests matrix operations to be computed by a set of worker processes. The system consists of three main components:

1. **Client:** Sends computation requests to the coordinator.
2. **Coordinator (Server):** Manages and assigns tasks to workers based on scheduling and load balancing.
3. **Workers:** Perform matrix operations and return results to the coordinator.

## Features
- Supports **Addition, Transpose, and Multiplication** matrix operations.
- Implements **Remote Procedure Calls (RPC)** for client-server communication.
- Uses **First-Come, First-Served (FCFS) Scheduling** for task management.
- Implements **Load Balancing** to assign tasks to the least busy worker.
- Provides **Fault Tolerance** by reassigning tasks if a worker fails.

## System Architecture
- The **client** interacts with the **coordinator** via RPC.
- The **coordinator** assigns matrix operations to available **worker** processes.
- The **worker** computes the requested matrix operation and sends the result back to the **coordinator**, which forwards it to the **client**.

## Requirements
- **Go** (Golang) for building the client, coordinator, and worker processes.
- **Docker** (optional) for containerized deployment.

## Project Structure
```
project_root/
├── client/
│   ├── go.mod
│   ├── main.go
│   ├── client.go
├── coordinator/
│   ├── go.mod
│   ├── main.go
│   ├── coordinator.go
├── worker/
│   ├── go.mod
│   ├── main.go
│   ├── worker.go
├── Dockerfile (optional)
└── README.md
```

## Installation and Setup
### 1. Clone the repository
```sh
git clone https://github.com/musaddiq-ashfaq/client-server-rpc
cd project_root
```

### 2. Build and Run the System
#### Running the Coordinator (Server)
```sh
cd coordinator
go run main.go
```

#### Running the Worker Processes
Each worker should be started separately:
```sh
cd worker
go run main.go
```

#### Running the Client
```sh
cd client
go run main.go
```

## Usage
- The **client** sends a matrix operation request (Addition, Transpose, Multiplication) to the **coordinator**.
- The **coordinator** assigns the task to the least busy **worker**.
- The **worker** computes the result and returns it to the **coordinator**.
- The **coordinator** sends the final result back to the **client**.

## Future Enhancements
- Implement secure communication using **TLS**.
- Improve fault tolerance using **checkpointing**.
- Implement **dynamic worker scaling** using Docker.

## Contributors
- Musaddiq Ashfaq
- Moaz farrukh
- Luqman Ansari


