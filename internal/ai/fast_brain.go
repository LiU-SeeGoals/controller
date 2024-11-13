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
		// fmt.Println("FastBrainGO: Sent actions")

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
		obstacleFreeDest := getObstacleFreeDest(robot, inst.Position, *gs)
		act.Dest = obstacleFreeDest

		act.Dribble = true // Assuming all moves require dribbling
		// fmt.Println("Team ", fb.team, ",Robot", act.Id, "moving:\n from", act.Pos.ToDTO(), "\n   to", act.Dest.ToDTO())
		// fmt.Println("Velocity: ", robot.GetVelocity())
		actionList = append(actionList, &act)
	}
	return actionList
}

const (
	fieldLength = 9000
	fieldWidth  = 6000
	cellSize    = 600
	kAtt        = 1.0 // Attractive potential constant
	kRep        = 2.0 // Repulsive potential constant
	d0          = 3.0 // Distance at which repulsive potential is 0
)

// robotPos is the current position of the robot
// goal is the destination position
// gs is the current game state, contains the positions of all robots ie. the obstacles
func getObstacleFreeDest(robot *state.Robot, goal state.Position, gs state.GameState) state.Position {

	// Matrix to hold the potential in the local neighborhood
	localGrid := mat.NewDense(20, 20, nil)
	fmt.Println(localGrid.At(0, 0))

	// obstacles := getObstacles(gs, robot.GetID())

	// manually set obstacle and robot for testing
	obstacle := state.Position{X: -2250, Y: 0, Z: 0, Angle: 0}
	obstacles := []state.Position{obstacle}
	x, y := float32(0.0), float32(0.0)
	goal = state.Position{X: -4500, Y: 0, Z: 0, Angle: 0}

	localGrid.Apply(func(i, j int, v float64) float64 {
		// x := robot.GetPosition().X/cellSize + float32(i-1)
		// y := robot.GetPosition().Y/cellSize + float32(j-1)
		x := x/cellSize + float32(i-7)
		y := y/cellSize + float32(j-7)

		// Compute the attractive potential
		attractive := computeAttractivePotential(x, y, goal.X/cellSize, goal.Y/cellSize)

		// Compute the repulsive potential
		repulsive := computeRepulsivePotential(x, y, obstacles, d0, kRep)

		return attractive + repulsive
	}, localGrid)
	
	// pretty print matrix
	fmt.Println("Local grid")
	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {
			fmt.Printf("%.2f ", localGrid.At(j, i))
		}
		fmt.Println()
	}
	fmt.Println("End local grid")

	minPotentialRow, minPotentialCol, _ := argmin(localGrid)

	// Calculate the offsets relative to the robot’s current position
	offsetX := float32(minPotentialRow-1) * cellSize
	offsetY := float32(minPotentialCol-1) * cellSize

	// Apply the offsets to the robot’s current position to get the new destination
	newX := robot.GetPosition().X + offsetX
	newY := robot.GetPosition().Y + offsetY

	return state.Position{X: newX, Y: newY}
}

func getObstacles(gs state.GameState, id state.ID) []state.Position {
	// Get all robots
	robots := append(gs.Blue_team[:], gs.Yellow_team[:]...)
	// Get the ball
	obstacles := make([]state.Position, len(robots)+1)
	for i, robot := range robots {
		if robot.GetID() != id {
			obstacles[i] = robot.GetPosition()
		}
	}
	return obstacles
}

func computeAttractivePotential(x, y, goalX, goalY float32) float64 {
	dx := math.Abs(float64(x - goalX))
	dy := math.Abs(float64(y - goalY))
	return 0.5 * kAtt * math.Sqrt(math.Pow(dx, 2)+math.Pow(dy, 2))
}

// calculateRepulsivePotential calculates the repulsive potential from obstacles
func computeRepulsivePotential(x, y float32, obstacles []state.Position, d0, kRep float64) float64 {
	repulsive := 0.0
	for _, obstacle := range obstacles {
		obstacleX, obstacleY := obstacle.X/cellSize, obstacle.Y/cellSize
		dx := math.Abs(float64(x - obstacleX))
		dy := math.Abs(float64(y - obstacleY))

		distance := math.Sqrt(dx*dx+dy*dy) * cellSize
		if distance < d0 {
			repulsive += 0.5 * kRep * math.Pow(1/distance-1/d0, 2)
		}
	}
	return repulsive
}

func argmin(m *mat.Dense) (int, int, float64) {
	minValue := math.MaxFloat64
	minRow, minCol := 0, 0

	rows, cols := m.Dims()
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			val := m.At(i, j)
			if val < minValue {
				minValue = val
				minRow, minCol = i, j
			}
		}
	}
	return minRow, minCol, minValue
}
