package main

import (
	"database/sql"
	"fmt"
	"time"
	t "time"
)

const (
	insertSnapshotQuery       string = `INSERT INTO snapshot (camera, timestamp, photosmall, photo) VALUES ($1, $2, $3, $4);`
	selectLatestDate          string = `SELECT MAX(timestamp) FROM snapshot;`
	selectDatesQuery          string = `SELECT DISTINCT timestamp FROM snapshot;`
	selectSnapshotAllCamQuery string = `SELECT camera, timestamp, photosmall, photo FROM snapshot, photosmall WHERE snapshot.id = photosmall.id AND timestamp >= $1;`
	selectSnapshotOneCamQuery string = `SELECT camera, timestamp, photosmall, photo FROM snapshot, photosmall WHERE snapshot.id = photosmall.id AND timestamp >= $1 AND camera = $2;`
)

func InsertSnapshot(camera string, timestamp t.Time, photosmall []byte, photo []byte) {
	psClient := GetPostgresqlClient()
	idPhotoSmall := InsertPhotoSmall(photosmall)
	idPhoto := InsertPhoto(photo)
	_, err := psClient.Db.Exec(insertSnapshotQuery, camera, timestamp, idPhotoSmall, idPhoto)
	if err != nil {
		panic(err)
	}
	fmt.Println("Snapshot inserted")
}

func GetLatestTimestampInserted() *time.Time {
	var timestamp time.Time
	psClient := GetPostgresqlClient()

	row := psClient.Db.QueryRow(selectLatestDate)
	switch err := row.Scan(&timestamp); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		return &timestamp
	default:
		panic(err)
	}
	return nil
}
