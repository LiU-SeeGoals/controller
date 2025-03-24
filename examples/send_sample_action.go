package examples

import (
	"math"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

func sendSampleAction() {
	BaseStationClient := client.NewBaseStationClient("127.0.0.1:25565")

	actions := []action.Action{
		&action.Stop{Id: 2},
		&action.MoveTo{
			Id:   3,
			Pos:  info.Position{X: 100, Y: 200, Z: 0, Angle: math.Pi}, // Example values for Pos
			Dest: info.Position{X: 200, Y: 300, Z: 0, Angle: math.Pi}, // Example values for Dest
		},
	}

	// Create a list of actions with different robot IDs

	BaseStationClient.SendActions(actions)
}
