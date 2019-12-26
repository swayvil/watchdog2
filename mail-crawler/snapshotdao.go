package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	LIMIT_SELECT              int    = 30
	insertSnapshotQuery       string = `INSERT INTO snapshot (camera, timestamp) VALUES ($1, $2);`
	selectLatestDate          string = `SELECT MAX(timestamp) FROM snapshot;`
	selectDatesQuery          string = `SELECT DISTINCT timestamp FROM snapshot;`
	countSnapshotAllCamQuery  string = `SELECT count(*) FROM snapshot WHERE timestamp >= $1;`
	selectSnapshotAllCamQuery string = `SELECT camera, timestamp FROM snapshot WHERE timestamp >= $1 LIMIT $2 OFFSET $3;`
	selectSnapshotOneCamQuery string = `SELECT camera, timestamp FROM snapshot WHERE timestamp >= $1 AND camera = $2 LIMIT $3 OFFSET $4;`
)

type Snapshot struct {
	Camera         string `json:"camera"`
	Timestamp      string `json:"timestamp"`
	PhotosmallPath string `json:"photosmallPath"`
	PhotoPath      string `json:"photoPath"`
}

func InsertSnapshot(camera string, timestamp time.Time, photosmall []byte, photo []byte) {
	if camera == "" {
		log.Fatal("InsertSnapshot: camera is null")
	}
	if photosmall == nil {
		log.Fatal("InsertSnapshot: photosmall is null")
	}
	if photo == nil {
		log.Fatal("InsertSnapshot: photo is null")
	}
	psClient := GetPostgresqlClient()
	fileName := timestamp.Format("20060102150405") + ".jpg"
	WriteImageToFs(photosmall, camera+string(os.PathSeparator)+"small"+string(os.PathSeparator)+fileName)
	WriteImageToFs(photo, camera+string(os.PathSeparator)+"big"+string(os.PathSeparator)+fileName)
	//photoSmallID := InsertPhotoSmall(photosmall)
	//photoID := InsertPhoto(photo)
	_, err := psClient.Db.Exec(insertSnapshotQuery, camera, timestamp)
	if err != nil {
		panic(err)
	}
	fmt.Println("Snapshot inserted")
}

func GetLatestTimestampInserted() *time.Time {
	var timestamp time.Time
	psClient := GetPostgresqlClient()

	row := psClient.Db.QueryRow(selectLatestDate)
	err := row.Scan(&timestamp)
	if err != nil {
		return nil
	}
	return &timestamp
}

func GetCountSnapshotAllCam(fromTimestamp time.Time) int {
	var count int
	psClient := GetPostgresqlClient()

	row := psClient.Db.QueryRow(countSnapshotAllCamQuery, fromTimestamp)
	err := row.Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

func GetSnapshotAllCam(fromTimestamp time.Time, cursor int) []Snapshot {
	psClient := GetPostgresqlClient()

	rows, err := psClient.Db.Query(selectSnapshotAllCamQuery, fromTimestamp, LIMIT_SELECT, cursor*LIMIT_SELECT)
	if err != nil {
		panic(err)
	}
	var snapshots []Snapshot
	snapshots = make([]Snapshot, LIMIT_SELECT)
	i := 0
	for rows.Next() {
		var camera string
		var timestamp time.Time
		err = rows.Scan(&camera, &timestamp)
		if err != nil {
			panic(err)
		}
		camera = strings.TrimSpace(camera)
		photosRoot := GetConfigInstance().WebServer.PhotosRoot
		fileName := timestamp.Format("20060102150405") + ".jpg"
		photosmallPath := photosRoot + "/" + camera + "/small/" + fileName
		photoPath := photosRoot + "/" + camera + "/big/" + fileName
		snapshot := Snapshot{camera, timestamp.Format("2006-01-02T15:04:05"), photosmallPath, photoPath}
		snapshots[i] = snapshot
		i++
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	if i < LIMIT_SELECT {
		return snapshots[0:i]
	} else {
		return snapshots
	}
}
