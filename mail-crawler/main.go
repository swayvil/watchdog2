package main

import (
	"fmt"
)

func main() {
	startMailCron()
	handleRequests()
	getPostgresqlClient().closeConnection()
	fmt.Println("Done")
}
