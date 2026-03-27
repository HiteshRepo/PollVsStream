package main

import grpcfibonacciclient "github.com/hiteshrepo/go-rest-vs-grpc/apps/grpc-fibonacci-client"

func main() {
	app := grpcfibonacciclient.NewApp()
	app.Start()

	app.Shutdown()
}
