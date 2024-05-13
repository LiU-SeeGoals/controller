package demos

import (
	"fmt"
	"time"

	"github.com/LiU-SeeGoals/controller/internal/ai"
	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/config"
	state "github.com/LiU-SeeGoals/controller/internal/gamestate"
	"github.com/LiU-SeeGoals/controller/internal/simulator"
)

func SimulatedAnnealing() {
	gs := state.NewGameState(3)
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

	sim_controller.RobotStartPositionConfig1(len(presentYellow))

	ssl_receiver.InitGameState(gs, 0)
	start_time := time.Now().UnixMilli()
	for {
		play_time := time.Now().UnixMilli() - start_time
		ssl_receiver.UpdateGamestate(gs, play_time)

		blue_actions, score_blue, antBlue := ai_blue.CreateActions(gs)
		yellow_actions, score_yellow, antYellow := ai_yellow.CreateActions(gs)

		sim_client_blue.SendActions(blue_actions)
		sim_client_yellow.SendActions(yellow_actions)

		terminal_messages := []string{fmt.Sprintf("Blue score: %.2f AnticipatedScore: %.2f", score_blue, antBlue), fmt.Sprintf("Yellow score: %.2f AnticipatedScore: %.2f", score_yellow, antYellow)}

		client.UpdateWebGUI(gs, blue_actions, terminal_messages)
		fmt.Printf("Blue/Yellow score: %.2f/%.2f, Anticipated score: %.2f/%.2f\n", score_blue, score_yellow, antBlue, antYellow)
	}
}
