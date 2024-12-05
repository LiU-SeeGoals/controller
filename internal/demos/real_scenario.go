package demos

import (
	"time"

	"github.com/LiU-SeeGoals/controller/internal/ai"
	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/config"
	"github.com/LiU-SeeGoals/controller/internal/state"
)

func RealScenario() {
	gs := state.NewGameState(10)
	ssl_receiver := client.NewSSLVisionClient(config.GetSSLClientAddressReal())

	// Yellow team
	slowBrainYellow := ai.NewScenarioSlowBrain(1, 2)
	fastBrainYellow := ai.NewFastBrainGO()

	aiYellow := ai.NewAi(state.Yellow, slowBrainYellow, fastBrainYellow)

	clientYellow := client.NewBaseStationClient(config.GetBasestationAddress())

	ssl_receiver.InitGameState(gs, 0)
	start_time := time.Now().UnixMilli()
	for {
		playTime := time.Now().UnixMilli() - start_time

		ssl_receiver.UpdateGamestate(gs, playTime)

		yellow_actions := aiYellow.GetActions(gs)
		clientYellow.SendActions(yellow_actions)

	}
}
