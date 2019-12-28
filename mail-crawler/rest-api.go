package main

import (
	"encoding/json"
	"fmt"
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
		panic(err)
	}

	cursor, err := strconv.Atoi(vars["cursor"])
	if err != nil {
		panic(err)
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
		panic(err)
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

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	photosRoot := getConfigInstance().WebServer.PhotosRoot
	myRouter.PathPrefix(photosRoot).Handler(http.StripPrefix(photosRoot, http.FileServer(http.Dir(getConfigInstance().Fs.PhotosStorePath))))
	myRouter.HandleFunc("/snapshots/{time}/{cams}/{cursor}", getSnapshotsAPI).Methods("GET")
	myRouter.HandleFunc("/count-snapshots/{time}/{cams}", countSnapshotsAPI).Methods("GET")
	myRouter.HandleFunc("/snapshots-limit", getSnapshotsLimitAPI).Methods("GET")
	myRouter.HandleFunc("/cameras", getCamerasAPI).Methods("GET")
	fmt.Println("Listening to 9999")

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST)
	handler := cors.Default().Handler(myRouter)
	log.Fatal(http.ListenAndServe(":9999", handler))
}
