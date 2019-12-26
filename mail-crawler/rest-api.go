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

func restSnapshotAllCam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ts := vars["ts"] // timestamp
	fromTimestamp, err := time.Parse(tsLayout, ts)
	cursor, err := strconv.Atoi(vars["cursor"])

	if err != nil {
		panic(err)
	}
	snapshots := GetSnapshotAllCam(fromTimestamp, cursor)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snapshots)
}

func restCountSnapshotAllCam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ts := vars["ts"] // timestamp
	fromTimestamp, err := time.Parse(tsLayout, ts)

	if err != nil {
		panic(err)
	}
	count := GetCountSnapshotAllCam(fromTimestamp)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(count)
}

func restSnapshoLimit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LIMIT_SELECT)
}

func HandleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	photosRoot := GetConfigInstance().WebServer.PhotosRoot
	myRouter.PathPrefix(photosRoot).Handler(http.StripPrefix(photosRoot, http.FileServer(http.Dir(GetConfigInstance().Fs.PhotosStorePath))))
	myRouter.HandleFunc("/snapshots-all-cams/{ts}/{cursor}", restSnapshotAllCam).Methods("GET")
	myRouter.HandleFunc("/count-snapshots-all-cams/{ts}", restCountSnapshotAllCam).Methods("GET")
	myRouter.HandleFunc("/snapshots-limit", restSnapshoLimit).Methods("GET")
	fmt.Println("Listening to 9999")

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST)
	handler := cors.Default().Handler(myRouter)
	log.Fatal(http.ListenAndServe(":9999", handler))
}
