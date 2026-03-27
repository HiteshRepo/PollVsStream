package main

import restfibonacciclient "github.com/hiteshrepo/go-rest-vs-grpc/apps/rest-fibonacci-client"

func main() {
	app := restfibonacciclient.NewApp()
	app.Start()
}
