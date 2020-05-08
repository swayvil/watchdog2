package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	insertSnapshotQuery   string = `INSERT INTO snapshot (camera, timestamp) VALUES ($1, $2);`
	selectLatestDateQuery string = `SELECT MAX(timestamp) FROM snapshot;`
	selectDatesQuery      string = `SELECT DISTINCT timestamp FROM snapshot;`
	countSnapshotsQuery   string = `SELECT count(*) FROM snapshot WHERE timestamp >= $1 AND camera = ANY($2::text[]);`
	selectSnapshotsQuery  string = `SELECT camera, timestamp FROM snapshot WHERE timestamp >= $1 AND camera = ANY($2::text[]) LIMIT $3 OFFSET $4;`
)

type snapshot struct {
	Camera         string `json:"camera"`
	Timestamp      string `json:"timestamp"`
	PhotosmallPath string `json:"photosmallPath"`
	PhotoPath      string `json:"photoPath"`
}

func insertSnapshot(camera string, timestamp time.Time, photosmall []byte, photo []byte) {
	if camera == "" {
		log.Fatal("InsertSnapshot: camera is null")
	}
	if photosmall == nil {
		log.Fatal("InsertSnapshot: photosmall is null")
	}
	if photo == nil {
		log.Fatal("InsertSnapshot: photo is null")
	}
	psClient := getPostgresqlClient()
	fileName := timestamp.Format("20060102150405") + ".jpg"
	writeImageToFs(photosmall, camera+string(os.PathSeparator)+"small"+string(os.PathSeparator)+fileName)
	writeImageToFs(photo, camera+string(os.PathSeparator)+"big"+string(os.PathSeparator)+fileName)
	//photoSmallID := InsertPhotoSmall(photosmall)
	//photoID := InsertPhoto(photo)
	_, err := psClient.Db.Exec(insertSnapshotQuery, camera, timestamp)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Snapshot inserted %s\n", timestamp)
}

func selectLatestTimestampInserted() *time.Time {
	var timestamp time.Time
	psClient := getPostgresqlClient()

	row := psClient.Db.QueryRow(selectLatestDateQuery)
	err := row.Scan(&timestamp)
	if err != nil {
		return nil
	}
	return &timestamp
}

func selectCountSnapshots(fromTimestamp time.Time, cameras string) int {
	var count int
	psClient := getPostgresqlClient()

	row := psClient.Db.QueryRow(countSnapshotsQuery, fromTimestamp, "{"+cameras+"}")
	err := row.Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

func selectSnapshots(fromTimestamp time.Time, cameras string, cursor int) []snapshot {
	limitSelect := getConfigInstance().Db.LimitSelect
	psClient := getPostgresqlClient()

	rows, err := psClient.Db.Query(selectSnapshotsQuery, fromTimestamp, "{"+cameras+"}", limitSelect, cursor*limitSelect)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	var snapshots []snapshot
	snapshots = make([]snapshot, limitSelect)
	i := 0
	for rows.Next() {
		var camera string
		var timestamp time.Time
		err = rows.Scan(&camera, &timestamp)
		if err != nil {
			panic(err)
		}
		camera = strings.TrimSpace(camera)
		photosRoot := getConfigInstance().WebServer.PhotosRoot
		fileName := timestamp.Format("20060102150405") + ".jpg"
		photosmallPath := photosRoot + "/" + camera + "/small/" + fileName
		photoPath := photosRoot + "/" + camera + "/big/" + fileName
		snapshot := snapshot{camera, timestamp.Format("2006-01-02T15:04:05"), photosmallPath, photoPath}
		snapshots[i] = snapshot
		i++
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	if i < limitSelect {
		return snapshots[0:i]
	}
	return snapshots
}
