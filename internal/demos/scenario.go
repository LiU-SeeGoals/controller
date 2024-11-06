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

func Scenario() {
	gs := state.NewGameState(10)
	ssl_receiver := client.NewSSLVisionClient(config.GetSSLClientAddress())

	slowBrainYellow := ai.NewScenarioSlowBrain(-5, -1)
	fastBrainYellow := ai.NewFastBrainGO()

	aiYellow := ai.NewAi(state.Yellow, slowBrainYellow, fastBrainYellow)

	simClientYellow := client.NewSimClient(config.GetSimYellowTeamAddress(), gs)

	simController := simulator.NewSimControl()

	// Some sim setup for debugging ai behaviour
	presentYellow := []int{0, 1}
	presentBlue := []int{}
	simController.SetPresentRobots(presentYellow, presentBlue)

	ssl_receiver.InitGameState(gs, 0)
	start_time := time.Now().UnixMilli()
	for {
		playTime := time.Now().UnixMilli() - start_time
		fmt.Println("playTime: ", playTime)
		ssl_receiver.UpdateGamestate(gs, playTime)

		yellow_actions := aiYellow.GetActions(gs)

		simClientYellow.SendActions(yellow_actions)

		// terminal_messages := []string{"Scenario"}

		// if len(blue_actions) > 0 {
		// 	client.UpdateWebGUI(gs, blue_actions, terminal_messages)
		// }
	}
}
