package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const (
	tsLayout string = "2006-01-02T15:04:05"
)

func getSnapshotsAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ts := vars["time"] + "T00:00:00" // timestamp
	fromTimestamp, err := time.Parse(tsLayout, ts)
	if err != nil {
		log.Fatal(err)
	}

	cursor, err := strconv.Atoi(vars["cursor"])
	if err != nil {
		log.Fatal(err)
	}
	cameras := vars["cams"]

	snapshots := selectSnapshots(fromTimestamp, cameras, cursor)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snapshots)
}

func countSnapshotsAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ts := vars["time"] + "T00:00:00" // timestamp
	fromTimestamp, err := time.Parse(tsLayout, ts)
	if err != nil {
		log.Fatal(err)
	}

	cameras := vars["cams"]
	count := selectCountSnapshots(fromTimestamp, cameras)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(count)
}

func getSnapshotsLimitAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getConfigInstance().Db.LimitSelect)
}

func getCamerasAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(selectCameras())
}

func getFirstSnapshotDate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(selectFirstTimestampInserted())
}

func getLastSnapshotDate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	timestamp := *selectLatestTimestampInserted()
	timestamp = timestamp.Truncate(24 * time.Hour) // Return only the date part
	json.NewEncoder(w).Encode(timestamp)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	photosRoot := getConfigInstance().WebServer.PhotosRoot
	listenPort := strconv.Itoa(getConfigInstance().WebServer.ListenPort)
	myRouter.PathPrefix(photosRoot).Handler(http.StripPrefix(photosRoot, http.FileServer(http.Dir(getConfigInstance().Fs.PhotosStorePath))))
	myRouter.HandleFunc("/snapshots/{time}/{cams}/{cursor}", getSnapshotsAPI).Methods("GET")
	myRouter.HandleFunc("/count-snapshots/{time}/{cams}", countSnapshotsAPI).Methods("GET")
	myRouter.HandleFunc("/snapshots-limit", getSnapshotsLimitAPI).Methods("GET")
	myRouter.HandleFunc("/cameras", getCamerasAPI).Methods("GET")
	myRouter.HandleFunc("/first-snapshot-date", getFirstSnapshotDate).Methods("GET")
	myRouter.HandleFunc("/last-snapshot-date", getLastSnapshotDate).Methods("GET")
	log.Println("Listening to " + listenPort)

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST)
	handler := cors.Default().Handler(myRouter)
	log.Fatal(http.ListenAndServe(":"+listenPort, handler))
}
