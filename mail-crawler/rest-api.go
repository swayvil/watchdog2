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

func returnSnapshotAllCam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ts := vars["ts"]
	fromTimestamp, err := time.Parse(tsLayout, ts)
	offset, err := strconv.Atoi(vars["os"])

	if err != nil {
		panic(err)
	}
	snapshots := GetSnapshotAllCam(fromTimestamp, offset)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snapshots)
}

func HandleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	photosRoot := GetConfigInstance().WebServer.PhotosRoot
	myRouter.PathPrefix(photosRoot).Handler(http.StripPrefix(photosRoot, http.FileServer(http.Dir(GetConfigInstance().Fs.PhotosStorePath))))
	myRouter.HandleFunc("/snapshots-all-cams/{ts}/{os}", returnSnapshotAllCam).Methods("GET")
	fmt.Println("Listening to 9999")

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST)
	handler := cors.Default().Handler(myRouter)
	log.Fatal(http.ListenAndServe(":9999", handler))
}
