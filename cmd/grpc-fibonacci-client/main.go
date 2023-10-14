package main

import grpcfibonacciclient "github.com/hiteshrepo/fibonacci-api/apps/grpc-fibonacci-client"

func main() {
	app := grpcfibonacciclient.NewApp()
	app.Start()

	app.Shutdown()
}
