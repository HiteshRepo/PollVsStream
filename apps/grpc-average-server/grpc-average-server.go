package grpcaverageserver

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/hiteshrepo/go-rest-vs-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	proto.UnimplementedAverageServiceServer
}

func NewApp() *App {
	return &App{}
}

func (a *App) Start() {
	serverPort := 50052
	servAddr := fmt.Sprintf("0.0.0.0:%d", serverPort)

	fmt.Println("starting average gRPC server at", servAddr)

	lis, err := net.Listen("tcp", servAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterAverageServiceServer(s, a)
	reflection.Register(s)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down gRPC average server...")
	s.GracefulStop()
}

func (a *App) Shutdown() {}

func (a *App) FindAverage(stream proto.AverageService_FindAverageServer) error {
	sum := 0
	count := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&proto.AverageResponse{
				Average: int32(sum / count),
			})
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
		}

		count += 1
		sum += int(req.GetNumber())
	}
}
