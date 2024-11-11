package ai

import (
	"fmt"
	"math"
	"time"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/state"
	"gonum.org/v1/gonum/mat"
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
		// dest := foo(act.Pos, inst.Position, *gs)
		act.Dest = inst.Position

		act.Dribble = true // Assuming all moves require dribbling
		// fmt.Println("Team ", fb.team, ",Robot", act.Id, "moving:\n from", act.Pos.ToDTO(), "\n   to", act.Dest.ToDTO())
		// fmt.Println("Velocity: ", robot.GetVelocity())
		actionList = append(actionList, &act)
	}
	return actionList
}

func foo(pos state.Position, goaltemp state.Position, gs state.GameState) {
	// Define grid dimensions
	gridWidth := int(math.Abs(float64(goaltemp.X)-float64(pos.X)))
	gridHeight := int(math.Abs(float64(goaltemp.Y)-float64(pos.Y)))

	// Define parameters
	// alpha := 50.0
	// beta := 50.0
	s := 15.0
	r := 2.0

	// Goal and obstacle coordinates
	goal := [2]float64{float64(goaltemp.X), float64(goaltemp.Y)}
	obstacle := [2]float64{3, 3}

	// Create X and Y arrays based on grid dimensions
	x := make([]float64, gridWidth)
	y := make([]float64, gridHeight)
	for i := 0; i < gridWidth; i++ {
		x[i] = float64(i)
	}
	for j := 0; j < gridHeight; j++ {
		y[j] = float64(j)
	}

	// Create meshgrid for X, Y
	X := mat.NewDense(gridHeight, gridWidth, nil)
	Y := mat.NewDense(gridHeight, gridWidth, nil)
	for i := 0; i < gridHeight; i++ {
		for j := 0; j < gridWidth; j++ {
			X.Set(i, j, x[j])
			Y.Set(i, j, y[i])
		}
	}

	// Initialize delx, dely matrices
	delx := mat.NewDense(gridHeight, gridWidth, nil)
	dely := mat.NewDense(gridHeight, gridWidth, nil)

	// Iterate over the grid
	for i := 0; i < gridHeight; i++ {
		for j := 0; j < gridWidth; j++ {
			xVal := X.At(i, j)
			yVal := Y.At(i, j)

			// Calculate distances
			dGoal := math.Sqrt(math.Pow(goal[0]-xVal, 2) + math.Pow(goal[1]-yVal, 2))
			dObstacle := math.Sqrt(math.Pow(obstacle[0]-xVal, 2) + math.Pow(obstacle[1]-yVal, 2))

			// Calculate angles
			thetaGoal := math.Atan2(goal[1]-yVal, goal[0]-xVal)
			thetaObstacle := math.Atan2(obstacle[1]-yVal, obstacle[0]-xVal)

			// Apply conditions to calculate delx, dely based on obstacle and goal distances
			if dObstacle < r {
				delx.Set(i, j, math.Copysign(1, math.Cos(thetaObstacle)))
				dely.Set(i, j, math.Copysign(1, math.Sin(thetaObstacle)))
			} else if dObstacle > r+s {
				delx.Set(i, j, 50*s*math.Cos(thetaObstacle))
				dely.Set(i, j, 50*s*math.Sin(thetaGoal))
			} else {
				delx.Set(i, j, -120*(s+r-dObstacle)*math.Cos(thetaObstacle))
				dely.Set(i, j, -120*(s+r-dObstacle)*math.Sin(thetaObstacle))
			}

			if dGoal < r+s {
				delxVal := delx.At(i, j)
				if delxVal != 0 {
					delx.Set(i, j, delxVal+(50*(dGoal-r)*math.Cos(thetaGoal)))
					dely.Set(i, j, dely.At(i, j)+(50*(dGoal-r)*math.Sin(thetaGoal)))
				} else {
					delx.Set(i, j, 50*(dGoal-r)*math.Cos(thetaGoal))
					dely.Set(i, j, 50*(dGoal-r)*math.Sin(thetaGoal))
				}
			} else {
				delxVal := delx.At(i, j)
				if delxVal != 0 {
					delx.Set(i, j, delxVal+50*s*math.Cos(thetaGoal))
					dely.Set(i, j, dely.At(i, j)+50*s*math.Sin(thetaGoal))
				} else {
					delx.Set(i, j, 50*s*math.Cos(thetaGoal))
					dely.Set(i, j, 50*s*math.Sin(thetaGoal))
				}
			}
		}
	}




	// Print final delx and dely matrices for debugging
	fmt.Println("delx matrix:")
	matPrint(delx)
	fmt.Println("\ndely matrix:")
	matPrint(dely)

	

	// You may use a visualization library in Go like gonum/plot to visualize these results.
}

// Helper function to print matrices for debugging
func matPrint(X mat.Matrix) {
	fa := mat.Formatted(X, mat.Prefix(""), mat.Excerpt(5))
	fmt.Printf("%v\n", fa)
}

