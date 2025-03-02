package demos

// import (
// 	"time"

// 	"github.com/LiU-SeeGoals/controller/internal/ai"
// 	"github.com/LiU-SeeGoals/controller/internal/client"
// 	"github.com/LiU-SeeGoals/controller/internal/config"
// 	"github.com/LiU-SeeGoals/controller/internal/info"
// )

// func RealScenario() {
// 	gameInfo := info.NewGameInfo(10)
// 	ssl_receiver := client.NewSSLClient()

// 	// Yellow team
// 	slowBrainYellow := ai.NewScenarioSlowBrain(1, 2)
// 	fastBrainYellow := ai.NewFastBrainGO()

// 	aiYellow := ai.NewAi(info.Yellow, slowBrainYellow, fastBrainYellow)

// 	clientYellow := client.NewBaseStationClient(config.GetBasestationAddress())

// 	start_time := time.Now().UnixMilli()
// 	for {
// 		playTime := time.Now().UnixMilli() - start_time

// 		ssl_receiver.UpdateState(gameInfo, playTime)

// 		yellow_actions := aiYellow.GetActions(gameInfo)
// 		clientYellow.SendActions(yellow_actions)

// 	}
// }

import (
	"time"

	"github.com/LiU-SeeGoals/controller/internal/ai"
	slow_brain "github.com/LiU-SeeGoals/controller/internal/ai/slow_brain"
	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/config"
	"github.com/LiU-SeeGoals/controller/internal/info"
	"github.com/LiU-SeeGoals/controller/internal/simulator"
)

func RealScenario() {
	// This avoid the "No position in history" error for robots
	presentYellow := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	presentBlue := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	simController := simulator.NewSimControl()
	simController.SetPresentRobots(presentYellow, presentBlue)

	gameInfo := info.NewGameInfo(10)
	ssl_receiver := client.NewSSLClient()

	// Yellow team
	//slowBrainYellow := slow_brain.NewSlowBrain1(info.Yellow)
	slowBrainYellow := slow_brain.NewSlowBrain2(info.Yellow)
	fastBrainYellow := ai.NewFastBrainGO()

	aiYellow := ai.NewAi(info.Yellow, slowBrainYellow, fastBrainYellow)

	simClientYellow := client.NewBaseStationClient(config.GetBasestationAddress())

	// Some sim setup for debugging ai behaviour
	presentYellow = []int{0, 1}
	presentBlue = []int{}
	simController.SetPresentRobots(presentYellow, presentBlue)

	start_time := time.Now().UnixMilli()
	for {
		playTime := time.Now().UnixMilli() - start_time
		ssl_receiver.UpdateState(gameInfo, playTime)
		yellow_actions := aiYellow.GetActions(gameInfo)
		simClientYellow.SendActions(yellow_actions)
	}
}
