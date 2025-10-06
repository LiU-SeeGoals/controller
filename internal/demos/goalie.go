package demos

import (
	"time"

	"github.com/LiU-SeeGoals/controller/internal/ai"
	slow_brain "github.com/LiU-SeeGoals/controller/internal/ai/slow_brain"
	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/config"
	"github.com/LiU-SeeGoals/controller/internal/info"
	"github.com/LiU-SeeGoals/controller/internal/simulator"
)

func Goalie() {
	// This avoid the "No position in history" error for robots
	presentYellow := []int{0}
	presentBlue := []int{0, 1}
	simController := simulator.NewSimControl()
	simController.SetPresentRobots(presentYellow, presentBlue)

	gameInfo := info.NewGameInfo(10)
	ssl_receiver := client.NewSSLClient(config.GetSSLClientAddress())

	// Yellow team
	slowBrainYellow := slow_brain.NewSlowBrainGoalie(info.Yellow)
	fastBrainYellow := ai.NewFastBrainGO()

	aiYellow := ai.NewAi(info.Yellow, slowBrainYellow, fastBrainYellow)

	simClientYellow := client.NewSimClient(config.GetSimYellowTeamAddress(), gameInfo)

	// Blue team
	slowBrainBlue := slow_brain.NewSlowBrainGoalie(info.Blue)
	fastBrainBlue := ai.NewFastBrainGO()

	aiBlue := ai.NewAi(info.Blue, slowBrainBlue, fastBrainBlue)

	simClientBlue := client.NewSimClient(config.GetSimBlueTeamAddress(), gameInfo)

	start_time := time.Now().UnixMilli()
	for {
		playTime := time.Now().UnixMilli() - start_time
		// fmt.Println("playTime: ", playTime)
		ssl_receiver.UpdateState(gameInfo, playTime)
		//fmt.Println(gameInfo)

		yellow_actions := aiYellow.GetActions(gameInfo)
		simClientYellow.SendActions(yellow_actions)

		blue_actions := aiBlue.GetActions(gameInfo)
		simClientBlue.SendActions(blue_actions)

		// if len(blue_actions) > 0 {
		// 	client.UpdateWebGUI(gameInfo, blue_actions, terminal_messages)
		// }
	}
}
