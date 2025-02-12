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
