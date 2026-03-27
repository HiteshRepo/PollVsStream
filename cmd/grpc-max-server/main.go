package main

import grpcmaxserver "github.com/hiteshrepo/go-rest-vs-grpc/apps/grpc-max-server"

func main() {
	app := grpcmaxserver.NewApp()
	app.Start()
}
