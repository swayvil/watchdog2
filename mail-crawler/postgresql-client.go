package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type postgresqlClient struct { //Singleton
	Database *sql.DB
}

var instanceDbClient *postgresqlClient = nil
var onceDbClient sync.Once
var connectAttempt int = 3
var waitInS time.Duration = 10

func getPostgresqlClient() *postgresqlClient {
	onceDbClient.Do(func() {
		var errPing error = nil

		connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			getConfigInstance().Db.Host, getConfigInstance().Db.Port, getConfigInstance().Db.User, getConfigInstance().Db.Password, getConfigInstance().Db.Name)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Panicln(err)
		}
		for connectAttempt > 0 {
			errPing = db.Ping()
			if errPing != nil {
				log.Printf("Will attempt %d more time to connect to Database\n", connectAttempt)
				log.Printf("Waiting %d s\n", waitInS)
				connectAttempt--
				time.Sleep(waitInS * time.Second)
			} else {
				instanceDbClient = &postgresqlClient{db}
				connectAttempt = 0 // We are connected
			}
		}
		if instanceDbClient == nil {
			if errPing != nil {
				log.Println(errPing)
			}
			log.Fatalln("Error while connecting to database, exiting")
		}
	})
	return instanceDbClient
}

func (psClient *postgresqlClient) closeConnection() {
	psClient.Database.Close()
	log.Printf("Database connection stopped")
}
