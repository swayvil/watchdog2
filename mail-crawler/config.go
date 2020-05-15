package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

const configFile string = "config.json"

type config struct {
	Imap struct {
		URL      string `json:"url"`
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"imap"`
	Db struct {
		Host        string `json:"host"`
		Port        int    `json:"port"`
		User        string `json:"user"`
		Password    string `json:"password"`
		Name        string `json:"name"`
		LimitSelect int    `json:"limitSelect"`
	} `json:"db"`
	Mail struct {
		Since struct {
			Year     int    `json:"year"`
			Month    int    `json:"month"`
			Day      int    `json:"day"`
			TimeZone string `json:"timeZone"`
		} `json:"since"`
		SubjectPattern string `json:"subjectPattern"`
		BodyPattern    string `json:"bodyPattern"`
	} `json:"mail"`
	Fs struct {
		PhotosStorePath string `json:"photosStorePath"`
	} `json:"fs"`
	WebServer struct {
		PhotosRoot string `json:"photosRoot"`
		ListenPort int    `json:"listenPort"`
	} `json:"web-server"`
}

var instanceConfig *config
var onceConfig sync.Once

func getConfigInstance() *config {
	onceConfig.Do(func() {
		instanceConfig = &config{}
		instanceConfig.loadConfiguration(configFile)
	})
	return instanceConfig
}

func (config *config) loadConfiguration(file string) {
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
}
