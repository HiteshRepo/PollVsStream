package main

import grpcfibonacciserver "github.com/hiteshrepo/fibonacci-api/apps/grpc-fibonacci-server"

func main() {
	app := grpcfibonacciserver.NewApp()
	app.Start()
}
