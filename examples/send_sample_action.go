package examples

import (
	"math"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/client"
	"gonum.org/v1/gonum/mat"
)

func sendSampleAction() {
	BaseStationClient := client.NewBaseStationClient("127.0.0.1:25565")

	actions := []action.Action{
		&action.Stop{Id: 2},
		&action.MoveTo{
			Id:   3,
			Pos:  mat.NewVecDense(3, []float64{100, 200, math.Pi}),  // Example values for Pos
			Dest: mat.NewVecDense(3, []float64{300, 400, -math.Pi}), // Example values for Goal
		},
	}

	// Create a list of actions with different robot IDs

	BaseStationClient.SendActions(actions)
}
