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
)

func RealScenario() {
	// This avoid the "No position in history" error for robots

	gameInfo := info.NewGameInfo(10)
	ssl_receiver := client.NewSSLClient(config.GetSSLClientAddressReal())

	// Yellow team
	//slowBrainYellow := slow_brain.NewSlowBrain1(info.Yellow)
	slowBrainYellow := slow_brain.NewSlowBrain1(info.Yellow)
	fastBrainYellow := ai.NewFastBrainGO()

	aiYellow := ai.NewAi(info.Yellow, slowBrainYellow, fastBrainYellow)
	//simClientYellow := client.Ne(config.GetSimYellowTeamAddress(), gameInfo)

	simClientYellow := client.NewBaseStationClient(config.GetBasestationAddress())

	simClientYellow.Init()

	start_time := time.Now().UnixMilli()
	start := time.Now()
	for {
		time.Sleep(500 * time.Millisecond)

		playTime := time.Now().UnixMilli() - start_time

		ssl_receiver.UpdateState(gameInfo, playTime)
		yellow_actions := aiYellow.GetActions(gameInfo)
		if time.Since(start) > 500 {
			simClientYellow.SendActions(yellow_actions)
			start = time.Now()
		}

	}
}
