package main

const selectCamerasQuery string = `SELECT camera FROM camera;`

func selectCameras() []string {
	psClient := getPostgresqlClient()
	rows, err := psClient.Db.Query(selectCamerasQuery)
	defer rows.Close()
	var cameras []string
	for rows.Next() {
		var camera string
		err = rows.Scan(&camera)
		if err != nil {
			// handle this error
			panic(err)
		}
		cameras = append(cameras, camera)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return cameras
}
