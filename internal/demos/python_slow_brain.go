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

func PythonSlowBrain() {
	gs := state.NewGameState(10)
	ssl_receiver := client.NewSSLVisionClient(config.GetSSLClientAddress())

	slowBrainBlue := ai.NewSlowBrainPy("http://10.240.211.153:5000/slowBrainBlue")
	fastBrainBlue := ai.NewFastBrainGO()

	aiBlue := ai.NewAi(state.Blue, slowBrainBlue, fastBrainBlue)

	slowBrainYellow := ai.NewScenarioSlowBrain(-5)
	fastBrainYellow := ai.NewFastBrainGO()

	aiYellow := ai.NewAi(state.Yellow, slowBrainYellow, fastBrainYellow)

	simClientBlue := client.NewSimClient(config.GetSimBlueTeamAddress())
	simClientYellow := client.NewSimClient(config.GetSimYellowTeamAddress())

	simController := simulator.NewSimControl()

	// Some sim setup for debugging ai behaviour
	presentYellow := []int{0, 1}
	presentBlue := []int{0, 1}
	simController.SetPresentRobots(presentYellow, presentBlue)

	ssl_receiver.InitGameState(gs, 0)
	start_time := time.Now().UnixMilli()
	for {
		playTime := time.Now().UnixMilli() - start_time
		fmt.Println("playTime: ", playTime)
		ssl_receiver.UpdateGamestate(gs, playTime)

		blue_actions := aiBlue.GetActions(gs)
		yellow_actions := aiYellow.GetActions(gs)

		simClientBlue.SendActions(blue_actions)
		simClientYellow.SendActions(yellow_actions)

		// terminal_messages := []string{"Scenario"}

		// if len(blue_actions) > 0 {
		// 	client.UpdateWebGUI(gs, blue_actions, terminal_messages)
		// }
	}
}
