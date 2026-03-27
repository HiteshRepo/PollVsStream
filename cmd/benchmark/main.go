package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/hiteshrepo/go-rest-vs-grpc/apps/benchmark"
)

func main() {
	n := flag.Int("n", 20, "number of fibonacci numbers to compute per run")
	iter := flag.Int("iter", 3, "number of benchmark iterations")
	pollMs := flag.Int("poll", 50, "REST polling interval in milliseconds")
	flag.Parse()

	restAddr := "localhost:8081"
	grpcAddr := "localhost:50055"

	fmt.Printf("Starting inline REST server on %s...\n", restAddr)
	stopREST := benchmark.StartRESTServer(restAddr)
	defer stopREST()

	fmt.Printf("Starting inline gRPC server on %s...\n", grpcAddr)
	stopGRPC := benchmark.StartGRPCServer(grpcAddr)
	defer stopGRPC()

	// warmup
	time.Sleep(100 * time.Millisecond)

	pollInterval := time.Duration(*pollMs) * time.Millisecond
	fmt.Printf("Running %d iterations (n=%d, poll=%v)...\n\n", *iter, *n, pollInterval)

	restResults := benchmark.BenchmarkREST(restAddr, *n, *iter, pollInterval)
	grpcResults := benchmark.BenchmarkGRPC(grpcAddr, *n, *iter)

	benchmark.PrintTable(restResults, grpcResults, *n)
}
