package examples

import (
	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/client"
	"gonum.org/v1/gonum/mat"
)

// Simple example showing how GrsimClient
//
// Example creates a client and sends 4 actions.
// Actions are ordered by id, every robot except id 3
// receives a stop action.
func GrsimClientExample() {
	c := client.NewGrsimClient("127.0.0.1:20011")
	c.Connect()

	move := &action.Move{
		Pos:     mat.NewVecDense(3, []float64{1.0, 1.0, 0.0}),
		Dest:    mat.NewVecDense(3, []float64{5.0, 1.0, 0.0}),
		Dribble: false,
	}

	stop := &action.Stop{}

	actions := []action.Action{stop, stop, move, stop}
	c.AddActions(actions)
}
