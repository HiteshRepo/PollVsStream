package main

import restfibonacciserver "github.com/hiteshrepo/go-rest-vs-grpc/apps/rest-fibonacci-server"

func main() {
	app := restfibonacciserver.NewApp()
	app.Start()
}
