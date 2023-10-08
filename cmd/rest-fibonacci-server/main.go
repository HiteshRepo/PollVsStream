package main

import restfibonacciserver "github.com/hiteshrepo/fibonacci-api/apps/rest-fibonacci-server"

func main() {
	app := restfibonacciserver.NewApp()
	app.Start()
}
