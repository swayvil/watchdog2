package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

type PostgresqlClient struct { //Singleton
	Db *sql.DB
}

var instanceDbClient *PostgresqlClient
var onceDbClient sync.Once

func GetPostgresqlClient() *PostgresqlClient {
	onceDbClient.Do(func() {
		connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			GetConfigInstance().Db.Host, GetConfigInstance().Db.Port, GetConfigInstance().Db.User, GetConfigInstance().Db.Password, GetConfigInstance().Db.Name)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal(err)
		}
		instanceDbClient = &PostgresqlClient{db}
	})
	return instanceDbClient
}

func (psClient *PostgresqlClient) CloseConnection() {
	psClient.Db.Close()
}
