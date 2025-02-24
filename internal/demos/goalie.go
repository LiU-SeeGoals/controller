package demos

import (
	"time"

	. "github.com/LiU-SeeGoals/controller/internal/logger"
	"github.com/LiU-SeeGoals/controller/internal/ai"
	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/config"
	"github.com/LiU-SeeGoals/controller/internal/info"
	"github.com/LiU-SeeGoals/controller/internal/simulator"
)

func Goalie() {
	gameInfo := info.NewGameInfo(10)
	ssl_receiver := client.NewSSLClient()

	slowBrainYellow := ai.NewScenarioSlowBrain(-5, 5)
	fastBrainYellow := ai.NewFastBrainGO()

	aiYellow := ai.NewAi(info.Yellow, slowBrainYellow, fastBrainYellow)

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
		// fmt.Println(gameInfo.Status)
		Logger.Info(gameInfo.Status)
		
		yellow_actions := aiYellow.GetActions(gameInfo)

		
		simClientYellow.SendActions(yellow_actions)

		// Communication to the GameViewer
		terminal_messages := []string{"Scenario Python Slow Brain"}
		client.UpdateWebGUI(gameInfo.State, yellow_actions, terminal_messages)
	}
}
