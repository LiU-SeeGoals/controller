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

func Scenario() {
	gameInfo := info.NewGameInfo(10)
	ssl_receiver := client.NewSSLClient()

	// Yellow team
	slowBrainYellow := ai.NewScenarioSlowBrain(1, 4)
	fastBrainYellow := ai.NewFastBrainGO()

	aiYellow := ai.NewAi(info.Yellow, slowBrainYellow, fastBrainYellow)

	simClientYellow := client.NewSimClient(config.GetSimYellowTeamAddress(), gameInfo)

	// Blue team
	//slowBrainBlue := ai.NewScenarioSlowBrain(1, 4)
	//fastBrainBlue := ai.NewFastBrainGO()

	//aiBlue := ai.NewAi(info.Blue, slowBrainBlue, fastBrainBlue)

	//simClientBlue := client.NewSimClient(config.GetSimBlueTeamAddress(), gameInfo)

	simController := simulator.NewSimControl()

	// Some sim setup for debugging ai behaviour
	presentYellow := []int{0, 1, 2, 3}
	presentBlue := []int{}
	simController.SetPresentRobots(presentYellow, presentBlue)

	start_time := time.Now().UnixMilli()
	for {
		playTime := time.Now().UnixMilli() - start_time
		// fmt.Println("playTime: ", playTime)
		ssl_receiver.UpdateState(gameInfo, playTime)
		fmt.Println(gameInfo)

		yellow_actions := aiYellow.GetActions(gameInfo)
		simClientYellow.SendActions(yellow_actions)

		//blue_actions := aiBlue.GetActions(gameInfo)
		//simClientBlue.SendActions(blue_actions)

		// terminal_messages := []string{"Scenario"}

		// if len(blue_actions) > 0 {
		// 	client.UpdateWebGUI(gameInfo, blue_actions, terminal_messages)
		// }
	}
}
