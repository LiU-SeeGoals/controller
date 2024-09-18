package main

import (
	"github.com/LiU-SeeGoals/controller/internal/ai"
	"github.com/LiU-SeeGoals/controller/internal/config"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"github.com/LiU-SeeGoals/controller/internal/receiver"
)

func main() {
	gs := gamestate.NewGameState()
	ssl_receiver := receiver.NewSSLReceiver(config.GetSSLClientAddress())
	ai := ai.NewAi(config.GetSimYellowTeamAddress(), gs)
	sim_controller := config.NewSimControl()

	presentYellow := []int{0, 1, 2, 3, 4, 5}
	presentBlue := []int{0, 1, 2, 3, 4, 5}
	sim_controller.SetPresentRobots(presentYellow, presentBlue)
	for {
		ssl_receiver.UpdateGamestate(gs)
		ai.CreateAndSendActions()
		// fmt.Println(*gs)
	}
}
