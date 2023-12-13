package main

import (
	"encoding/json"
	"os"

	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"github.com/LiU-SeeGoals/controller/internal/webserver"
)

type Config struct {
	SSLClientAddress string `json:"sslClientAddress"`
	GrSimAddress     string `json:"grSimAddress"`
}

func main() {
	conf, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	go webserver.Once.Do(webserver.StartWebServer)


	gs := gamestate.NewGameState(conf.GrSimAddress, conf.SSLClientAddress)
	for {
		gs.TestActions()
		gs.Update()

		//fmt.Println(gs)
	}
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("../config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
