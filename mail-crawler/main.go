package main

import (
	"fmt"
)

func main() {
	//dao.InsertSnapshot("garage", time.Now(), make([]byte, 5), make([]byte, 5))
	imapCLient := NewImapClient()
	imapCLient.GetMessages()
	imapCLient.Logout()
	GetPostgresqlClient().CloseConnection()
	fmt.Println("Done")
}
