package demos

import (
	"fmt"
	"time"

	"github.com/LiU-SeeGoals/controller/internal/ai"
	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/config"
	"github.com/LiU-SeeGoals/controller/internal/simulator"
	"github.com/LiU-SeeGoals/controller/internal/state"
)

func GoTo() {
	gs := state.NewGameState(10)
	ssl_receiver := client.NewSSLVisionClient(config.GetSSLClientAddress())

	ai_blue := ai.NewAi(state.Blue)
	sim_client_blue := client.NewSimClient(config.GetSimBlueTeamAddress())

	sim_controller := simulator.NewSimControl()

	// Some sim setup for debugging ai behaviour
	presentYellow := []int{1, 4}
	presentBlue := []int{2}
	sim_controller.SetPresentRobots(presentYellow, presentBlue)

	ssl_receiver.InitGameState(gs, 0)
	start_time := time.Now().UnixMilli()
	for {
		play_time := time.Now().UnixMilli() - start_time
		fmt.Println("play_time: ", play_time)
		ssl_receiver.UpdateGamestate(gs, play_time)

		blue_actions := ai_blue.GetActions(gs)

		sim_client_blue.SendActions(blue_actions)

		terminal_messages := []string{"GoTo"}

		if len(blue_actions) > 0 {
			client.UpdateWebGUI(gs, blue_actions, terminal_messages)
		}
	}
}
