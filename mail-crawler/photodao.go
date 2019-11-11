package main

import (
	"log"
)

const (
	insertPhotoQuery string = `INSERT INTO photo (photo) VALUES ($1) RETURNING id;`
	selectPhotoQuery string = `SELECT photo FROM photo WHERE photo.id = $1;`
)

func InsertPhoto(photo []byte) int {
	psClient := GetPostgresqlClient()

	id := 0
	err := psClient.Db.QueryRow(insertPhotoQuery, photo).Scan(&id)
	if err != nil {
		log.Println("Insert photo failed.")
	}
	return id
}
