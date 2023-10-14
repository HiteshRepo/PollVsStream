package main

import grpcmaxclient "github.com/hiteshrepo/fibonacci-api/apps/grpc-max-client"

func main() {
	app := grpcmaxclient.NewApp()
	app.Start()

	app.Shutdown()
}
