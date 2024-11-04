# REST Polling vs gRPC Streaming

This repository provides a comparison between REST-based polling and gRPC-based streaming for real-time data delivery. It explores the strengths, weaknesses, and use cases for each approach.

## Overview

REST-based polling and gRPC-based streaming are two common methods for fetching real-time data in distributed systems. This project demonstrates their implementations and compares their performance and resource utilization under similar conditions.

## Features

- **REST Polling:** A REST API endpoint that provides data updates through periodic polling.
- **gRPC Streaming:** A gRPC server that streams data to clients, providing continuous updates.

## Project Structure

- `apps/`: Contains the server and client implementations of REST and gRPC supporting polling/streaming.
- `cmd/`: Contains starting points for those server and client implementations.
- `proto/`: Message and RPC definitions of the gRPC server-client contracts.

## Prerequisites

- Go (or the language of your choice)

## Usage

1. **Clone the repository:**
   ```bash
   git clone git@github.com:HiteshRepo/PollVsStream.git
   cd PollVsStream
   ```
2. **For REST**
    ```bash
    # on terminal-1
    go run cmd/rest-fibonacci-server/main.go

    # on terminal-2
    go run cmd/rest-fibonacci-client/main.go
    ```
3. **For gRPC client-side streaming**
    ```bash
    # on terminal-1
    go run cmd/grpc-average-server/main.go

    # on terminal-2
    go run cmd/grpc-average-client/main.go
    ```
4. **For gRPC server-side streaming**
    ```bash
    # on terminal-1
    go run cmd/grpc-fibonacci-server/main.go

    # on terminal-2
    go run cmd/grpc-fibonacci-client/main.go
    ```
5. **For gRPC server-side streaming**
    ```bash
    # on terminal-1
    go run cmd/grpc-max-server/main.go

    # on terminal-2
    go run cmd/grpc-max-client/main.go
    ```

## TODOs

1. Dockerize all the apps & group relevant ones into a single docker-compose.
2. Come up with K8s manifests.
3. Add README for each relevant groups.