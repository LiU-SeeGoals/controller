package main

import (
	// "fmt"

	"time"

	"github.com/LiU-SeeGoals/controller/internal/ai"
	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/config"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"github.com/LiU-SeeGoals/controller/internal/receiver"
	"github.com/LiU-SeeGoals/controller/internal/simulator"
	"github.com/LiU-SeeGoals/controller/internal/webserver"
)

func main() {
	start_time := time.Now().UnixMilli()
	gs := gamestate.NewGameState(3)
	ssl_receiver := receiver.NewSSLReceiver(config.GetSSLClientAddress())

	ai_blue := ai.NewAi(gamestate.Blue)
	sim_client_blue := client.NewSimClient(config.GetSimBlueTeamAddress())

	ai_yellow := ai.NewAi(gamestate.Yellow)
	sim_client_yellow := client.NewSimClient(config.GetSimYellowTeamAddress())

	sim_controller := simulator.NewSimControl()

	presentYellow := []int{0, 1, 2, 3, 4, 5}
	presentBlue := []int{0, 1, 2, 3, 4, 5}
	sim_controller.SetPresentRobots(presentYellow, presentBlue)

	terminal_messages := []string{"message 1", "message 2"}
	ssl_receiver.InitGameState(gs, time.Now().UnixMilli()-start_time)
	for {
		play_time := time.Now().UnixMilli() - start_time
		ssl_receiver.UpdateGamestate(gs, play_time)

		blue_actions := ai_blue.CreateActions(gs)
		yellow_actions := ai_yellow.CreateActions(gs)

		sim_client_blue.SendActions(blue_actions)
		sim_client_yellow.SendActions(yellow_actions)

		webserver.UpdateWebGUI(gs, blue_actions, terminal_messages)
		// fmt.Println(*gs)
	}
}
