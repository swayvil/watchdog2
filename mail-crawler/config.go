package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

const configFile string = "config.json"

type Config struct {
	Imap struct {
		Url      string `json:"url"`
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"imap"`
	Db struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Name     string `json:"name"`
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
}

var instanceConfig *Config
var onceConfig sync.Once

func GetConfigInstance() *Config {
	onceConfig.Do(func() {
		instanceConfig = &Config{}
		instanceConfig.loadConfiguration(configFile)
	})
	return instanceConfig
}

func (config *Config) loadConfiguration(file string) {
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
}
