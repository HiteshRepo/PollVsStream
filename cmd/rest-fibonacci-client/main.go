package main

import restfibonacciclient "github.com/hiteshrepo/fibonacci-api/apps/rest-fibonacci-client"

func main() {
	app := restfibonacciclient.NewApp()
	app.Start()
}
