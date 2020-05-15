package main

import (
	"log"
	"strings"
)

const selectCamerasQuery string = `SELECT camera FROM camera;`

func selectCameras() []string {
	psClient := getPostgresqlClient()
	rows, err := psClient.Database.Query(selectCamerasQuery)
	defer rows.Close()
	var cameras []string
	for rows.Next() {
		var camera string
		err = rows.Scan(&camera)
		if err != nil {
			log.Fatal(err)
		}
		cameras = append(cameras, strings.TrimSpace(camera))
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	if cameras == nil || len(cameras) <= 0 {
		log.Fatalln("No cameras found in database, exiting")
	}
	return cameras
}
