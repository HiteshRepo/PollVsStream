# REST Polling vs gRPC Streaming

This repository provides a comparison between REST-based polling and gRPC-based streaming for real-time data delivery. It explores the strengths, weaknesses, and use cases for each approach.

## Overview

REST-based polling and gRPC-based streaming are two common methods for fetching real-time data in distributed systems. This project demonstrates their implementations and compares their performance and resource utilization under similar conditions.

## Features

- **REST Polling:** A REST API endpoint that provides data updates through periodic polling.
- **gRPC Streaming:** A gRPC server that streams data to clients, providing continuous updates.
- **Graceful Shutdown:** All servers handle `SIGINT`/`SIGTERM` cleanly.
- **Benchmark Tool:** Self-contained benchmark that spins up its own servers and compares REST polling vs gRPC streaming side-by-side.
- **Docker Support:** Multi-stage Dockerfile and `docker-compose.yml` for all servers.

## Project Structure

```
apps/               # Server and client implementations
  benchmark/        # Benchmark engine (inline REST + gRPC servers)
  grpc-*/           # gRPC server and client packages
  rest-fibonacci-*/ # REST server and client packages
cmd/                # Entry points
  benchmark/        # Benchmark CLI
  grpc-*/
  rest-fibonacci-*/
proto/              # Protobuf definitions and generated Go code
```

## Prerequisites

- Go 1.20+
- Docker & Docker Compose (optional)
- `protoc` + `protoc-gen-go` + `protoc-gen-go-grpc` (only if regenerating proto files)

## Usage

### Run individual apps

```bash
# REST (port 8080)
go run cmd/rest-fibonacci-server/main.go
go run cmd/rest-fibonacci-client/main.go -typeOfCall async -number 20

# gRPC server-side streaming — fibonacci (port 50051)
go run cmd/grpc-fibonacci-server/main.go
go run cmd/grpc-fibonacci-client/main.go -typeOfCall async -number 20

# gRPC client-side streaming — average (port 50052)
go run cmd/grpc-average-server/main.go
go run cmd/grpc-average-client/main.go -numbers "1,2,3,4,5"

# gRPC bidirectional streaming — max (port 50053)
go run cmd/grpc-max-server/main.go
go run cmd/grpc-max-client/main.go -numbers "3,1,4,1,5,9"
```

### Run the benchmark

The benchmark starts its own inline servers — no external processes needed.

```bash
go run cmd/benchmark/main.go -n 20 -iter 5 -poll 50
```

| Flag | Default | Description |
|------|---------|-------------|
| `-n` | `20` | Number of Fibonacci numbers per run |
| `-iter` | `3` | Number of iterations |
| `-poll` | `50` | REST polling interval (ms) |

### Make targets

```bash
make test       # go test ./...
make benchmark  # run the benchmark tool
make docker-up  # docker compose up --build
make gen-proto  # regenerate protobuf Go files
```

### Docker

```bash
docker compose up --build
```

Starts all four servers:

| Service | Port |
|---------|------|
| rest-fibonacci-server | 8080 |
| grpc-fibonacci-server | 50051 |
| grpc-average-server | 50052 |
| grpc-max-server | 50053 |

## TODOs

1. Add K8s manifests.
2. Add per-service READMEs.
