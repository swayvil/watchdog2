package main

import (
	"fmt"
)

func main() {
	//imapCLient := NewImapClient()
	//imapCLient.ImportMessages()
	//imapCLient.Logout()
	HandleRequests()
	GetPostgresqlClient().CloseConnection()
	fmt.Println("Done")
}
