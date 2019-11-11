package main

import (
	"log"
)

const (
	insertPhotoSmallQuery string = `INSERT INTO photosmall (photo) VALUES ($1) RETURNING id;`
	selectPhotoSmallQuery string = `SELECT photo FROM photosmall WHERE photo.id = $1;`
)

func InsertPhotoSmall(photo []byte) int {
	psClient := GetPostgresqlClient()

	id := 0
	err := psClient.Db.QueryRow(insertPhotoSmallQuery, photo).Scan(&id)
	if err != nil {
		log.Println("Insert photo failed.")
	}
	return id
}
