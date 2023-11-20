package main

import (
	"math"
	"time"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/client"
	"gonum.org/v1/gonum/mat"
)


func main() {
	BaseStationClient := client.NewBaseStationClient("127.0.0.1:25565")
	BaseStationClient.Init()

	actions := []action.Action{
		&action.Stop{Id: 2},
		&action.Move{
			Id: 3,
			Pos: mat.NewVecDense(3, []float64{100, 200, math.Pi}), // Example values for Pos
			Goal: mat.NewVecDense(3, []float64{300, 400, -math.Pi}), // Example values for Goal
		},
	}

		// Create a list of actions with different robot IDs

	BaseStationClient.Send(actions)

	time.Sleep(2 * time.Second)
	BaseStationClient.Send(actions)
	time.Sleep(2 * time.Second)
}