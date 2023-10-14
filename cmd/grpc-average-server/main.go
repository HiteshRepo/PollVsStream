package main

import grpcaverageserver "github.com/hiteshrepo/fibonacci-api/apps/grpc-average-server"

func main() {
	app := grpcaverageserver.NewApp()
	app.Start()
}
