package benchmark

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/hiteshrepo/go-rest-vs-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// Result holds timing data for one benchmark run.
type Result struct {
	FirstLatency time.Duration
	TotalTime    time.Duration
	RoundTrips   int
}

// ---- inline fib ----

func inlineFib(n int) int {
	if n <= 0 {
		return 0
	} else if n == 1 {
		return 1
	}
	return inlineFib(n-1) + inlineFib(n-2)
}

// ---- inline REST server ----

type inlineAsyncStore struct {
	numbers        []int
	current        int
	requestedRange int
}

func newInlineAsyncStore(r int) *inlineAsyncStore { return &inlineAsyncStore{requestedRange: r} }

func (s *inlineAsyncStore) write(number, current int) {
	s.numbers = append(s.numbers, number)
	s.current = current
}

func (s *inlineAsyncStore) read() ([]int, int, int) {
	nums := make([]int, len(s.numbers))
	copy(nums, s.numbers)
	s.numbers = nil
	return nums, s.current, s.requestedRange
}

type inlineRESTApp struct {
	stores map[string]*inlineAsyncStore
}

func (a *inlineRESTApp) asyncHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n := 0
	fmt.Sscanf(vars["number"], "%d", &n)

	reqID := r.Header.Get("request-id")
	if reqID == "" {
		reqID = uuid.New().String()
	}

	if _, ok := a.stores[reqID]; !ok {
		a.stores[reqID] = newInlineAsyncStore(n)
		go func(store *inlineAsyncStore, count int) {
			for i := 0; i < count; i++ {
				store.write(inlineFib(i), i)
			}
		}(a.stores[reqID], n)
	}

	nums, current, requested := a.stores[reqID].read()
	end := (current+1) == requested

	type asyncResp struct {
		RequestID        string `json:"requestid"`
		FibonacciNumbers []int  `json:"fibonacciNumbers"`
		EndOfResponse    bool   `json:"endOfResponse"`
	}
	resp := asyncResp{RequestID: reqID, FibonacciNumbers: nums, EndOfResponse: end}
	if end {
		delete(a.stores, reqID)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// StartRESTServer starts an inline REST server on addr and returns a shutdown func.
func StartRESTServer(addr string) func() {
	app := &inlineRESTApp{stores: make(map[string]*inlineAsyncStore)}
	router := mux.NewRouter()
	router.HandleFunc("/fibonacci/async/{number:[0-9]+}", app.asyncHandler)

	srv := &http.Server{Addr: addr, Handler: router}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("REST benchmark server error: %v", err)
		}
	}()
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}
}

// ---- inline gRPC server ----

type inlineGRPCServer struct {
	proto.UnimplementedFibonacciServiceServer
}

func (s *inlineGRPCServer) AsyncFibonacci(req *proto.FibonacciRequest, stream proto.FibonacciService_AsyncFibonacciServer) error {
	n := int(req.GetNumber())
	for i := 0; i < n; i++ {
		if err := stream.Send(&proto.AsyncFibonacciResponse{
			Sequence:        int32(i),
			FibonacciNumber: int32(inlineFib(i)),
		}); err != nil {
			return err
		}
	}
	return nil
}

// StartGRPCServer starts an inline gRPC fibonacci server on addr and returns a shutdown func.
func StartGRPCServer(addr string) func() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("benchmark gRPC listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterFibonacciServiceServer(s, &inlineGRPCServer{})
	reflection.Register(s)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Printf("gRPC benchmark server error: %v", err)
		}
	}()
	return func() { s.GracefulStop() }
}

// ---- REST benchmark ----

// BenchmarkREST runs iter iterations of async REST polling for n fibonacci numbers.
func BenchmarkREST(restAddr string, n, iter int, pollInterval time.Duration) []Result {
	results := make([]Result, 0, iter)

	for i := 0; i < iter; i++ {
		url := fmt.Sprintf("http://%s/fibonacci/async/%d", restAddr, n)
		headers := map[string]string{"content-type": "application/json"}

		start := time.Now()
		var firstLatency time.Duration
		trips := 0
		first := true

		for {
			trips++
			reqID, ok := headers["request-id"]
			_ = reqID

			req, _ := http.NewRequest(http.MethodGet, url, nil)
			for k, v := range headers {
				req.Header.Set(k, v)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Printf("REST poll error: %v", err)
				break
			}

			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			if first {
				firstLatency = time.Since(start)
				first = false
			}

			type pollResp struct {
				RequestID     string `json:"requestid"`
				EndOfResponse bool   `json:"endOfResponse"`
			}
			var pr pollResp
			json.Unmarshal(body, &pr)

			if !ok && pr.RequestID != "" {
				headers["request-id"] = pr.RequestID
			}

			if pr.EndOfResponse {
				break
			}
			time.Sleep(pollInterval)
		}

		results = append(results, Result{
			FirstLatency: firstLatency,
			TotalTime:    time.Since(start),
			RoundTrips:   trips,
		})

		// reset request-id for next iteration
		delete(headers, "request-id")
	}

	return results
}

// ---- gRPC benchmark ----

// BenchmarkGRPC runs iter iterations of gRPC server-side streaming for n fibonacci numbers.
func BenchmarkGRPC(grpcAddr string, n, iter int) []Result {
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("benchmark gRPC dial: %v", err)
	}
	defer conn.Close()

	client := proto.NewFibonacciServiceClient(conn)
	results := make([]Result, 0, iter)

	for i := 0; i < iter; i++ {
		start := time.Now()
		var firstLatency time.Duration
		trips := 0
		first := true

		stream, err := client.AsyncFibonacci(context.Background(), &proto.FibonacciRequest{Number: int32(n)})
		if err != nil {
			log.Printf("gRPC stream error: %v", err)
			continue
		}

		for {
			_, err := stream.Recv()
			trips++
			if first {
				firstLatency = time.Since(start)
				first = false
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("gRPC recv error: %v", err)
				break
			}
		}

		results = append(results, Result{
			FirstLatency: firstLatency,
			TotalTime:    time.Since(start),
			RoundTrips:   trips,
		})
	}

	return results
}

// PrintTable prints an ASCII comparison table for REST vs gRPC results.
func PrintTable(restResults, grpcResults []Result, n int) {
	fmt.Println()
	fmt.Printf("┌─────────────────────────────────────────────────────────────────────────────────────────┐\n")
	fmt.Printf("│  Benchmark: REST async polling vs gRPC streaming  (n=%d fibonacci numbers)             │\n", n)
	fmt.Printf("├──────┬──────────────────────────────────────────┬──────────────────────────────────────┤\n")
	fmt.Printf("│      │              REST polling                │            gRPC streaming            │\n")
	fmt.Printf("│ Iter │ First Latency │  Total Time  │  Trips   │ First Latency │  Total Time  │ Msgs  │\n")
	fmt.Printf("├──────┼───────────────┼──────────────┼──────────┼───────────────┼──────────────┼───────┤\n")

	var restTotalFirst, restTotalTime time.Duration
	var grpcTotalFirst, grpcTotalTime time.Duration
	var restTotalTrips, grpcTotalTrips int

	for i := 0; i < len(restResults) && i < len(grpcResults); i++ {
		r := restResults[i]
		g := grpcResults[i]
		fmt.Printf("│ %4d │ %13s │ %12s │ %8d │ %13s │ %12s │ %5d │\n",
			i+1,
			r.FirstLatency.Round(time.Microsecond),
			r.TotalTime.Round(time.Microsecond),
			r.RoundTrips,
			g.FirstLatency.Round(time.Microsecond),
			g.TotalTime.Round(time.Microsecond),
			g.RoundTrips,
		)
		restTotalFirst += r.FirstLatency
		restTotalTime += r.TotalTime
		restTotalTrips += r.RoundTrips
		grpcTotalFirst += g.FirstLatency
		grpcTotalTime += g.TotalTime
		grpcTotalTrips += g.RoundTrips
	}

	count := time.Duration(len(restResults))
	if count == 0 {
		count = 1
	}
	fmt.Printf("├──────┼───────────────┼──────────────┼──────────┼───────────────┼──────────────┼───────┤\n")
	fmt.Printf("│  AVG │ %13s │ %12s │ %8d │ %13s │ %12s │ %5d │\n",
		(restTotalFirst / count).Round(time.Microsecond),
		(restTotalTime / count).Round(time.Microsecond),
		restTotalTrips/int(count),
		(grpcTotalFirst / count).Round(time.Microsecond),
		(grpcTotalTime / count).Round(time.Microsecond),
		grpcTotalTrips/int(count),
	)
	fmt.Printf("└──────┴───────────────┴──────────────┴──────────┴───────────────┴──────────────┴───────┘\n")
	fmt.Println()
}
