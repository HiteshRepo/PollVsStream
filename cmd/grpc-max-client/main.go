package main

import grpcmaxclient "github.com/hiteshrepo/go-rest-vs-grpc/apps/grpc-max-client"

func main() {
	app := grpcmaxclient.NewApp()
	app.Start()

	app.Shutdown()
}
