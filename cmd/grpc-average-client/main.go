package main

import grpcaverageclient "github.com/hiteshrepo/go-rest-vs-grpc/apps/grpc-average-client"

func main() {
	app := grpcaverageclient.NewApp()
	app.Start()

	app.Shutdown()
}
