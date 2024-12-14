package demos

import (
	"fmt"
	"time"

	"github.com/LiU-SeeGoals/controller/internal/ai"
	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/config"
	"github.com/LiU-SeeGoals/controller/internal/state"
)

func RealScenario() {
	gs := state.NewGameState(10)
	ssl_receiver := client.NewSSLClient(config.GetSSLClientAddressReal(), config.GetGCClientAddress())

	// Yellow team
	slowBrainYellow := ai.NewScenarioSlowBrain(1, 2)
	fastBrainYellow := ai.NewFastBrainGO()

	aiYellow := ai.NewAi(state.Yellow, slowBrainYellow, fastBrainYellow)
	fmt.Printf("address %s", config.GetBasestationAddress())
	clientYellow := client.NewBaseStationClient("10.242.33.22:6001")
	clientYellow.Init()
	ssl_receiver.InitState(gs, 0)
	start_time := time.Now().UnixMilli()
	for {

		playTime := time.Now().UnixMilli() - start_time

		ssl_receiver.UpdateState(gs, playTime)

		yellow_actions := aiYellow.GetActions(gs)
		if len(yellow_actions) > 0 {

			fmt.Println(yellow_actions[0])
		}
		time.Sleep(100 * time.Millisecond)
		clientYellow.SendActions(yellow_actions)

	}
}
