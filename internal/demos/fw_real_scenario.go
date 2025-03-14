package demos

import (
	"fmt"
	"time"

	"github.com/LiU-SeeGoals/controller/internal/ai"
	slow_brain "github.com/LiU-SeeGoals/controller/internal/ai/slow_brain"
	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/config"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

func FwRealScenario() {
	gameInfo := info.NewGameInfo(10)
	ssl_receiver := client.NewSSLClient(config.GetSSLClientAddressReal())

	// Yellow team
	slowBrainYellow := slow_brain.NewSlowBrainFw(info.Yellow)
	fastBrainYellow := ai.NewFastBrainGO()

	aiYellow := ai.NewAi(info.Yellow, slowBrainYellow, fastBrainYellow)

	basestationClient := client.NewBaseStationClient(config.GetBasestationAddress())
    fmt.Println("Base(d)station: ", config.GetBasestationAddress())

	basestationClient.Init()

	for {
		playTime := time.Now().UnixMilli()

		ssl_receiver.UpdateState(gameInfo, playTime)
		yellow_actions := aiYellow.GetActions(gameInfo)

		basestationClient.SendActions(yellow_actions)
	}
}
