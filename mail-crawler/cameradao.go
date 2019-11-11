package main

import (
	"fmt"
)

const selectCamerasQuery string = `SELECT camera FROM camera;`

func SelectCameras() {
	psClient := GetPostgresqlClient()
	rows, err := psClient.Db.Query(selectCamerasQuery)
	defer rows.Close()
	for rows.Next() {
		var camera string
		err = rows.Scan(&camera)
		if err != nil {
			// handle this error
			panic(err)
		}
		fmt.Println(camera)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}
}
