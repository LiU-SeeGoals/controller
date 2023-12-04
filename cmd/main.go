package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/LiU-SeeGoals/controller/internal/ai"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"github.com/LiU-SeeGoals/controller/internal/world_predictor"
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
	gs := gamestate.NewGameState()
	worldPredictor := world_predictor.NewWorldPredictor(conf.SSLClientAddress, gs)
	ai := ai.NewAi(gs, conf.GrSimAddress)

	for {
		ai.Update()
		worldPredictor.Update()
		fmt.Println(gs)
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
