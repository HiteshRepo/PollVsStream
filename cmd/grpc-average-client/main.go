package main

import grpcaverageclient "github.com/hiteshrepo/fibonacci-api/apps/grpc-average-client"

func main() {
	app := grpcaverageclient.NewApp()
	app.Start()

	app.Shutdown()
}
