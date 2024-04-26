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
	sim_config := config.NewSimControl()

	presentYellow := []int{0, 1, 2, 3, 4, 5, 6}
	presentBlue := []int{0, 1, 2}
	sim_config.SetPresentRobots(presentYellow, presentBlue)
	for {
		worldPredictor.Update()
		ai.Update()
		fmt.Println(gs)
	}
}
