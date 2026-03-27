package main

import grpcaverageserver "github.com/hiteshrepo/go-rest-vs-grpc/apps/grpc-average-server"

func main() {
	app := grpcaverageserver.NewApp()
	app.Start()
}
