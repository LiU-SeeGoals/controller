package demos

import (
	"fmt"
	"time"

	"github.com/LiU-SeeGoals/controller/internal/ai"
	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/config"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"github.com/LiU-SeeGoals/controller/internal/simulator"
	"github.com/LiU-SeeGoals/proto_go/simulation"
)

func GaussEnemy() {
	gs := gamestate.NewGameState(3)
	ssl_receiver := client.NewSSLVisionClient(config.GetSSLClientAddress())

	ai_blue := ai.NewAi(gamestate.Blue)
	sim_client_blue := client.NewSimClient(config.GetSimBlueTeamAddress())

	ai_yellow := ai.NewAi(gamestate.Yellow)
	sim_client_yellow := client.NewSimClient(config.GetSimYellowTeamAddress())

	sim_controller := simulator.NewSimControl()

	// Some sim setup for debugging ai behaviour
	presentYellow := []int{0, 1}
	presentBlue := []int{0}
	sim_controller.SetPresentRobots(presentYellow, presentBlue)
	sim_controller.TeleportRobot(0.1, 0.5, 0, simulation.Team_YELLOW)
	sim_controller.TeleportRobot(-0.5, -0.5, 1, simulation.Team_YELLOW)
	sim_controller.TeleportRobot(0, 0, 0, simulation.Team_BLUE)

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
	}
}