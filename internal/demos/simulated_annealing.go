package demos

import (
	"time"

	"github.com/LiU-SeeGoals/controller/internal/ai"
	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/config"
	"github.com/LiU-SeeGoals/controller/internal/simulator"
	"github.com/LiU-SeeGoals/controller/internal/state"
)

func SimulatedAnnealing() {
	gs := state.NewGameState(10)
	ssl_receiver := client.NewSSLVisionClient(config.GetSSLClientAddress())

	ai_blue := ai.NewAi(state.Blue)
	sim_client_blue := client.NewSimClient(config.GetSimBlueTeamAddress())

	ai_yellow := ai.NewAi(state.Yellow)
	sim_client_yellow := client.NewSimClient(config.GetSimYellowTeamAddress())

	sim_controller := simulator.NewSimControl()

	// Some sim setup for debugging ai behaviour
	presentYellow := []int{0, 1, 2, 3, 4, 5}
	presentBlue := []int{0, 1, 2, 3, 4, 5}
	sim_controller.SetPresentRobots(presentYellow, presentBlue)

	ssl_receiver.InitGameState(gs, 0)
	start_time := time.Now().UnixMilli()
	for {
		play_time := time.Now().UnixMilli() - start_time
		ssl_receiver.UpdateGamestate(gs, play_time)

		blue_actions := ai_blue.GetActions(gs)
		yellow_actions := ai_yellow.GetActions(gs)

		sim_client_blue.SendActions(blue_actions)
		sim_client_yellow.SendActions(yellow_actions)

		terminal_messages := []string{"Simulated Annealing"}

		client.UpdateWebGUI(gs, blue_actions, terminal_messages)
	}
}
