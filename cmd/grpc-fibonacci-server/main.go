package main

import grpcfibonacciserver "github.com/hiteshrepo/go-rest-vs-grpc/apps/grpc-fibonacci-server"

func main() {
	app := grpcfibonacciserver.NewApp()
	app.Start()
}
