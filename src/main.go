package main

import (
	"fmt"
	"os"

	"data"
	"fetch"
	"server"

	"encoding/json"
)

type Config struct {
	Server server.Config
	Db     data.Config
}

func main() {
	config := LoadConfiguration("config.json")

	fetch.InitFetcher()
	data.InitDB(config.Db)
	server.Start(config.Server)
}

func LoadConfiguration(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
