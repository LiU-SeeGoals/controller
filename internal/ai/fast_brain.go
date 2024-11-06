package ai

import (
	"fmt"
	"time"
	"gonum.org/v1/gonum/mat"
	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/state"
)

type FastBrainGO struct {
	team              state.Team
	incomingGameState <-chan state.GameState
	incomingGamePlan  <-chan state.GamePlan
	outgoingActions   chan<- []action.Action
}

func NewFastBrainGO() *FastBrainGO {
	return &FastBrainGO{}
}

func (fb *FastBrainGO) Init(incomingGameState <-chan state.GameState, incomingGamePlan <-chan state.GamePlan, outgoingActions chan<- []action.Action, team state.Team) {

	fb.incomingGameState = incomingGameState
	fb.incomingGamePlan = incomingGamePlan
	fb.outgoingActions = outgoingActions
	fb.team = team
	//

	go fb.Run()
}

func (fb *FastBrainGO) Run() {
	gameState := state.GameState{}
	gamePlan := state.GamePlan{}

	for {
		// We will reive the game state more often than the game plan
		// so we wait for the gameState to update and work with the latest game plan

		gameState = <-fb.incomingGameState

		select {
		case gamePlan = <-fb.incomingGamePlan:
		default:

		}
		// time.Sleep(1 * time.Second) // TODO: Remove this

		// Wait for the game to start
		if !gameState.Valid || !gamePlan.Valid {
			fmt.Println("FastBrainGO: Invalid game state")
			fb.outgoingActions <- []action.Action{}
			time.Sleep(10 * time.Millisecond)
			continue
		}

		// Do some thinking
		actions := fb.GetActions(&gameState, &gamePlan)

		// Send the actions to the AI
		fb.outgoingActions <- actions
		fmt.Println("FastBrainGO: Sent actions")

	}
}

func (fb *FastBrainGO) GetActions(gs *state.GameState, gamePlan *state.GamePlan) []action.Action {

	var actionList []action.Action

	myTeam := gs.GetTeam(fb.team)

	if fb.team != gamePlan.Team {
		panic("FastBrainGO: Team mismatch")
	}

	Instructions := gamePlan.Instructions

	for _, inst := range Instructions {
		robot := myTeam[inst.Id]

		if !robot.IsActive() {
			continue
		}
		act := action.MoveTo{}
		act.Id = int(inst.Id)
		act.Team = fb.team

		act.Pos = robot.GetPosition()

		act.Dest = inst.Position

		act.Dribble = true // Assuming all moves require dribbling
		// fmt.Println("Team ", fb.team, ",Robot", act.Id, "moving:\n from", act.Pos.ToDTO(), "\n   to", act.Dest.ToDTO())
		fmt.Println("Velocity: ", robot.GetVelocity())
		actionList = append(actionList, &act)
	}

	return actionList
}


func foo() {
	X := mat.NewDense(3, 3, []float64{0, 0, 0, 0, 0, 0, 0, 0, 0})
	Y := mat.NewDense(3, 3, []float64{0, 0, 0, 0, 0, 0, 0, 0, 0})
	
	s := 7.0
	r := 2.0

	// Creating two evenly spaced arrays ranging from -10 to 10
	var x, y []float64
	for i := -10; i < 10; i++ {
		x = append(x, float64(i))
		y = append(y, float64(i))
	}

	// Creating the meshgrid
	size := len(x)
	X := make([][]float64, size)
	Y := make([][]float64, size)
	delx := make([][]float64, size)
	dely := make([][]float64, size)

	for i := range X {
		X[i] = make([]float64, size)
		Y[i] = make([]float64, size)
		delx[i] = make([]float64, size)
		dely[i] = make([]float64, size)
		for j := range X[i] {
			X[i][j] = x[j]
			Y[i][j] = y[i]
			delx[i][j] = 0
			dely[i][j] = 0
		}
	}

	// Filling delx and dely arrays based on the conditions
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			d := math.Sqrt(X[i][j]*X[i][j] + Y[i][j]*Y[i][j])
			theta := math.Atan2(Y[i][j], X[i][j])

			if d < r {
				delx[i][j] = 0
				dely[i][j] = 0
			} else if d > r+s {
				delx[i][j] = -50 * s * math.Cos(theta)
				dely[i][j] = -50 * s * math.Sin(theta)
			} else {
				delx[i][j] = -50 * (d - r) * math.Cos(theta)
				dely[i][j] = -50 * (d - r) * math.Sin(theta)
			}
		}
	}



}

