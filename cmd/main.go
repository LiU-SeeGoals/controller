package main

import (
	"fmt"
	"time"

	"github.com/LiU-SeeGoals/controller/internal/ai"
	"github.com/LiU-SeeGoals/controller/internal/config"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"github.com/LiU-SeeGoals/controller/internal/world_predictor"
	"github.com/LiU-SeeGoals/proto_go/simulation"
)

func main() {
	// reader := bufio.NewReader(os.Stdin)
	gs := gamestate.NewGameState()
	worldPredictor := world_predictor.NewWorldPredictor(config.GetSSLClientAddress(), gs)
	sim_config := config.NewSimControl()
	ai := ai.NewAi(config.GetSimYellowTeamAddress(), gs, sim_config)

	presentYellow := []int{0, 1}
	presentBlue := []int{0}
	sim_config.SetPresentRobots(presentYellow, presentBlue)
	sim_config.TeleportRobot(0.1, 0.5, 0, simulation.Team_YELLOW)
	sim_config.TeleportRobot(-0.5, -0.5, 1, simulation.Team_YELLOW)
	sim_config.TeleportRobot(0, 0, 0, simulation.Team_BLUE)
	time.Sleep(1000 * time.Millisecond)
	for {
		worldPredictor.UpdateGamestate()
		ai.CreateAndSendActions()
		fmt.Println(gs.Blue_team[0].String())
		// _, _ = reader.ReadString('\n')
	}
}
