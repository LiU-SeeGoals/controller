package main

import (
	"fmt"

	"github.com/LiU-SeeGoals/controller/internal/ai"
	"github.com/LiU-SeeGoals/controller/internal/config"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"github.com/LiU-SeeGoals/controller/internal/world_predictor"
)

func main() {
	gs := gamestate.NewGameState(config.GetSSLClientAddress())
	worldPredictor := world_predictor.NewWorldPredictor(config.GetSSLClientAddress(), gs)
	ai := ai.NewAi(gs, config.GetSimYellowTeamAddress())
	for {
		worldPredictor.Update()
		ai.Update()
		fmt.Println(gs)
	}
}
