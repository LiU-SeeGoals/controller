package demos

import (
	"fmt"
	"time"

	"github.com/LiU-SeeGoals/controller/internal/ai"
	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/config"
	"github.com/LiU-SeeGoals/controller/internal/info"
	"github.com/LiU-SeeGoals/controller/internal/simulator"
)

func PythonSlowBrain() {
	gameInfo := info.NewGameInfo(10)
	ssl_receiver := client.NewSSLClient()

	slowBrainBlue := ai.NewSlowBrainPy("http://10.240.211.153:5000/slowBrainBlue")
	fastBrainBlue := ai.NewFastBrainGO()

	aiBlue := ai.NewAi(info.Blue, slowBrainBlue, fastBrainBlue)

	slowBrainYellow := ai.NewScenarioSlowBrain(-5, -1)
	fastBrainYellow := ai.NewFastBrainGO()

	aiYellow := ai.NewAi(info.Yellow, slowBrainYellow, fastBrainYellow)

	simClientBlue := client.NewSimClient(config.GetSimBlueTeamAddress(), gameInfo)
	simClientYellow := client.NewSimClient(config.GetSimYellowTeamAddress(), gameInfo)

	simController := simulator.NewSimControl()

	// Some sim setup for debugging ai behaviour
	presentYellow := []int{0, 1}
	presentBlue := []int{}
	simController.SetPresentRobots(presentYellow, presentBlue)

	start_time := time.Now().UnixMilli()
	for {
		playTime := time.Now().UnixMilli() - start_time
		// fmt.Println("playTime: ", playTime)
		ssl_receiver.UpdateState(gameInfo, playTime)
		fmt.Println(gameInfo.Status)
		blue_actions := aiBlue.GetActions(gameInfo)
		yellow_actions := aiYellow.GetActions(gameInfo)

		simClientBlue.SendActions(blue_actions)
		simClientYellow.SendActions(yellow_actions)

		// terminal_messages := []string{"Scenario"}

		// if len(blue_actions) > 0 {
		// 	client.UpdateWebGUI(gi, blue_actions, terminal_messages)
		// }
	}
}
