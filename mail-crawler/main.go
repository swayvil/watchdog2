package main

import "fmt"

func main() {
	//imapCLient := newImapClient()
	//imapCLient.importMessages()
	//imapCLient.logout()
	handleRequests()
	getPostgresqlClient().closeConnection()
	fmt.Println("Done")
}
