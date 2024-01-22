package main

import (
	"github.com/LiU-SeeGoals/controller/internal/ai"
	"github.com/LiU-SeeGoals/controller/internal/config"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"github.com/LiU-SeeGoals/controller/internal/world_predictor"
)

func main() {

	gs := gamestate.NewGameState()
	worldPredictor := world_predictor.NewWorldPredictor(config.GetSSLClientAddress(), gs)
	ai := ai.NewAi(gs, config.GetGrSimAddress())
	for {

		worldPredictor.Update()
		ai.Update()
	}

}
