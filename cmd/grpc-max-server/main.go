package main

import grpcmaxserver "github.com/hiteshrepo/fibonacci-api/apps/grpc-max-server"

func main() {
	app := grpcmaxserver.NewApp()
	app.Start()
}
