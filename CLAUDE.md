# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run any server or client
go run cmd/<app-name>/main.go

# Build all packages
go build ./...

# Run tests
go test ./...

# Regenerate protobuf Go files from .proto definitions
make gen-proto
# requires: protoc + protoc-gen-go + protoc-gen-go-grpc
```

## Architecture

This repo compares **REST polling** vs **gRPC streaming** for real-time data using Fibonacci numbers, average calculation, and max finding as workloads.

### Layout

- `proto/` — Protobuf definitions and generated Go code (`.pb.go`, `_grpc.pb.go`). Edit `.proto` files; regenerate with `make gen-proto`.
- `apps/<name>/` — Business logic packages (servers and clients).
- `cmd/<name>/main.go` — Entrypoints that instantiate and start the app from `apps/`.

### The three gRPC streaming patterns

| Service | Proto | Streaming type | Port |
|---|---|---|---|
| `FibonacciService` | `fibonacci.proto` | Unary + **server-side streaming** | 50051 |
| `AverageService` | `average.proto` | **Client-side streaming** | 50052 |
| `MaxService` | `max.proto` | **Bi-directional streaming** | 50053 |

### REST polling pattern (port 8080)

The REST server exposes two endpoints:
- `GET /fibonacci/sync/{n}` — blocks until all `n` Fibonacci numbers are computed, returns them at once.
- `GET /fibonacci/async/{n}` — starts a background goroutine keyed by `request-id` header, stores partial results in an in-memory `AsyncStore` (mutex-guarded). Clients poll repeatedly; the server returns accumulated numbers and sets `endOfResponse: true` on the final poll.

The middleware (`generateRequestIDMiddleware`) auto-generates a UUID `request-id` header if the client omits it. This ID is the key for async state.

### Module name

`github.com/hiteshrepo/go-rest-vs-grpc` — use this when adding new internal imports.
