package main

import (
	"log/slog"

	"github.com/LiU-SeeGoals/controller/internal/config"
	"github.com/LiU-SeeGoals/controller/internal/world_predictor"
)

func main() {
	config.SetLoggerConfig()
	slog.Info("Starting controller")
	worldPredictor := world_predictor.NewWorldPredictor()
	worldPredictor.Update()
	//ai := ai.NewAi(gs, config.GetGrSimAddress())

	for {
		//worldPredictor.Update()
		//ai.Update()
	}

}
