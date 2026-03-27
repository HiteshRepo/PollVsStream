package grpcmaxserver

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hiteshrepo/go-rest-vs-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	proto.UnimplementedMaxServiceServer
}

func NewApp() *App {
	return &App{}
}

func (a *App) Start() {
	serverPort := 50053
	servAddr := fmt.Sprintf("0.0.0.0:%d", serverPort)

	fmt.Println("starting max gRPC server at", servAddr)

	lis, err := net.Listen("tcp", servAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterMaxServiceServer(s, a)
	reflection.Register(s)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down gRPC max server...")
	s.GracefulStop()
}

func (a *App) Shutdown() {}

func (a *App) FindMax(stream proto.MaxService_FindMaxServer) error {
	max := int32(-1)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
		}

		num := req.GetNumber()
		if max < num {
			max = num
		}

		time.Sleep(2000 * time.Millisecond)
		err = stream.Send(&proto.MaxResponse{Max: max})
		if err != nil {
			log.Fatalf("Error while sending to client: %v", err)
			return err
		}
	}
}
