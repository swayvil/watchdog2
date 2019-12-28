package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

type postgresqlClient struct { //Singleton
	Db *sql.DB
}

var instanceDbClient *postgresqlClient
var onceDbClient sync.Once

func getPostgresqlClient() *postgresqlClient {
	onceDbClient.Do(func() {
		connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			getConfigInstance().Db.Host, getConfigInstance().Db.Port, getConfigInstance().Db.User, getConfigInstance().Db.Password, getConfigInstance().Db.Name)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal(err)
		}
		instanceDbClient = &postgresqlClient{db}
	})
	return instanceDbClient
}

func (psClient *postgresqlClient) closeConnection() {
	psClient.Db.Close()
}
