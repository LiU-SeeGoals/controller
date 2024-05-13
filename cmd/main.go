package main

import (
	"github.com/LiU-SeeGoals/controller/internal/ai"
	"github.com/LiU-SeeGoals/controller/internal/config"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"github.com/LiU-SeeGoals/controller/internal/world_predictor"
	"fmt"
)

func main() {
	gs := gamestate.NewGameState()
	worldPredictor := world_predictor.NewWorldPredictor(config.GetSSLClientAddress(), gs)
	ai := ai.NewAi(config.GetSimYellowTeamAddress(), gs)
	sim_config := config.NewSimControl()

	presentYellow := []int{0, 1, 2, 3, 4, 5}
	presentBlue := []int{0, 1, 2, 3, 4, 5}
	sim_config.SetPresentRobots(presentYellow, presentBlue)
	sim_config.RobotStartPositionConfig1(6)
	fmt.Println("TESTING HEHEHHEHE")
	for {
		worldPredictor.Update()
		ai.Update()
		// fmt.Println(*gs)
	}
}
