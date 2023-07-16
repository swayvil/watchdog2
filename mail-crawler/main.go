package main

func main() {
	defer getPostgresqlClient().closeConnection() // Close the DB connection when app stops
	createFolders()                               // Create local folders, if they don't exist, to store the snapshots
	startMailCron()                               // Schedule import every 30 minutes (for new snapshots)

	go getImapClient().importMessages() // Start the import
	handleRequests()                    // Start web server
	//getImapClient().deleteMessages("inbox", 2023, time.January, 1)
}
