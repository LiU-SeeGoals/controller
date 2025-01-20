package ai

import (
	"fmt"
	"math"
	"time"

	"github.com/LiU-SeeGoals/controller/internal/info"
)

const (
	RUNNING int = iota
	COMPLETE
	TIME_EXPIRED
	ERROR
	FAILED
)

type ScenarioSlowBrain struct {
	team              info.Team
	incomingGameState <-chan info.GameState
	outgoingPlan      chan<- info.GamePlan
	scale             float32
	run_scenario      int // -1 for all
}

func NewScenarioSlowBrain(scale float32, run_scenario int) *ScenarioSlowBrain {
	return &ScenarioSlowBrain{scale: scale, run_scenario: run_scenario}
}

func (sb *ScenarioSlowBrain) Init(incoming <-chan info.GameState, outgoing chan<- info.GamePlan, team info.Team) {
	sb.incomingGameState = incoming
	sb.outgoingPlan = outgoing
	sb.team = team

	go sb.Run()
}

type ScenarioTest interface {
	Run() []*info.Instruction
	Archived(*info.GameState) int
}

func (sb ScenarioSlowBrain) Run() {
	var gameState info.GameState
	gameState.SetValid(false)

	scenarios := []ScenarioTest{}
	scenarios = append(scenarios, NewMoveToTest(sb.team))
	scenarios = append(scenarios, NewMoveToBallTest(sb.team))
	scenarios = append(scenarios, NewObstacleAvoidanceTest(sb.team))
	scenarios = append(scenarios, NewRealTest(sb.team))
	// scenarios = append(scenarios, NewObstacleAvoidanceTest(sb.team))

	scenario_index := 0
	if sb.run_scenario >= 0 {
		scenario_index = sb.run_scenario
	}

	fmt.Println("Running scenarios")
	for {
		gameState = <-sb.incomingGameState

		if !gameState.IsValid() {
			fmt.Println("ScenarioSlowBrain: Invalid game state")
			time.Sleep(40 * time.Millisecond)
			continue
		}

		scenario := scenarios[scenario_index]
		game_state := scenario.Archived(&gameState)
		if game_state == COMPLETE {
			fmt.Println("Scenario", scenario_index, "completed")
		} else if game_state == TIME_EXPIRED {
			fmt.Println("Scenario", scenario_index, "time expired")
		} else if game_state == FAILED {
			fmt.Println("Scenario", scenario_index, "failed")
		}
		if game_state != RUNNING {
			scenario_index++
			if scenario_index >= len(scenarios) || sb.run_scenario >= 0 {
				panic("ScenarioSlowBrain: No more scenarios") // TODO: Handle this better
			}
			scenario = scenarios[scenario_index]
		}

		inst := scenario.Run()

		plan := info.GamePlan{}
		plan.Instructions = inst

		plan.Team = sb.team

		plan.Valid = true

		sb.outgoingPlan <- plan

	}
}

type MoveToBallTest struct {
	team     info.Team
	at_state int
	start    time.Time
	max_time time.Duration
}

func NewMoveToBallTest(team info.Team) *MoveToBallTest {
	return &MoveToBallTest{
		team:     team,
		max_time: 30 * time.Second,
		at_state: -1,
	}
}

func (m *MoveToBallTest) Run() []*info.Instruction {
	if m.at_state == -1 {
		m.start = time.Now()
		m.at_state = 0
	}
	if m.at_state == 0 {
		return []*info.Instruction{
			{Type: info.MoveToBall, Id: 0},
		}
	} else {
		return []*info.Instruction{
			{Type: info.MoveToBall, Id: 0},
		}
	}
}

func (m *MoveToBallTest) Archived(gs *info.GameState) int {
	robot_pos := gs.GetRobot(info.ID(0), m.team).GetPosition()
	ball_pos := gs.GetBall().GetPosition()

	dx := float64(robot_pos.X - ball_pos.X)
	dy := float64(robot_pos.Y - ball_pos.Y)
	distance := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))

	if m.at_state == 0 {
		if distance < 500 {
			m.at_state = 1
		}
	} else if m.at_state == 1 {
		if distance > 500 {
			fmt.Println("Failed with robot at (", robot_pos.X, robot_pos.Y, ") and ball at (", ball_pos.X, ball_pos.Y, ")")
			m.at_state = 2
		}
	}

	if m.at_state >= 0 {
		if time.Since(m.start) > m.max_time || m.at_state == 2 {
			if m.at_state == 0 {
				fmt.Println("Did not reach ball")
				return TIME_EXPIRED
			} else if m.at_state == 1 {
				fmt.Println("Reached ball and stayed there! :D")
				return COMPLETE
			} else {
				fmt.Println("Reached ball but then lost it :(")
				return FAILED
			}
		}
	}
	return RUNNING
}

// --------------------------Real Test--------------------------------
type RealTest struct {
	team           info.Team
	at_state       int
	start          time.Time
	max_time       time.Duration
	instructionSet []*info.Instruction
}

func NewRealTest(team info.Team) *RealTest {

	instSet := []*info.Instruction{
		{Type: info.MoveToPosition, Id: 5, Position: info.Position{X: -2000, Y: 0}},
		{Type: info.MoveToPosition, Id: 5, Position: info.Position{X: 2000, Y: 0}},
	}

	return &RealTest{
		team:           team,
		max_time:       600 * time.Second,
		at_state:       0,
		instructionSet: instSet,
	}
}

func (m *RealTest) Run() []*info.Instruction {

	return []*info.Instruction{m.instructionSet[m.at_state]}
}

func (m *RealTest) Archived(gs *info.GameState) int {
	// if len(gs.Yellow_team) != 4 {
	// 	fmt.Println("Yellow team has", len(gs.Yellow_team), "4 required")
	// 	return ERROR
	// } else if len(gs.Blue_team) != 3 {
	// 	fmt.Println("Blue team has", len(gs.Blue_team), "3 required")
	// 	return ERROR
	// }

	// This assumes that the robot ids range from 0 to m.team_size
	target := m.instructionSet[m.at_state].Position
	robot := gs.GetRobot(info.ID(5), m.team)

	if atPosition(robot, target) {
		m.at_state = (m.at_state + 1) % 2 // Cycle through the states
	}

	return RUNNING
}

// --------------------------Obstacle avoidance--------------------------------
type ObstacleAvoidanceTest struct {
	team           info.Team
	team_size      int
	at_states      []int
	start          time.Time
	max_time       time.Duration
	instructionSet [][]info.Instruction
}

func NewObstacleAvoidanceTest(team info.Team) *ObstacleAvoidanceTest {
	team_size := 4
	var instSet [][]info.Instruction
	if team == info.Yellow {
		instSet = [][]info.Instruction{
			{
				{Type: info.MoveToPosition, Id: 0, Position: info.Position{X: -4000, Y: 0}},
				{Type: info.MoveToPosition, Id: 0, Position: info.Position{X: 4000, Y: 0}},
			},
			{
				{Type: info.MoveToPosition, Id: 1, Position: info.Position{X: 640, Y: -1000}},
				{Type: info.MoveToPosition, Id: 1, Position: info.Position{X: 640, Y: 1000}},
			},
			{ // --------------------------Obstacle avoidance--------------------------------
				{Type: info.MoveToPosition, Id: 2, Position: info.Position{X: 1925, Y: -1000}},
				{Type: info.MoveToPosition, Id: 2, Position: info.Position{X: 1925, Y: 1000}},
			},
			{

				{Type: info.MoveToPosition, Id: 3, Position: info.Position{X: 3210, Y: -1000}},
				{Type: info.MoveToPosition, Id: 3, Position: info.Position{X: 3210, Y: 1000}},
			},
			{

				{Type: info.MoveToPosition, Id: 3, Position: info.Position{X: 3210, Y: -1000}},
				{Type: info.MoveToPosition, Id: 3, Position: info.Position{X: 3210, Y: 1000}},
			},
		}
	} else {
		instSet = [][]info.Instruction{
			{
				{Type: info.MoveToPosition, Id: 0, Position: info.Position{X: -3215, Y: 1000}},
				{Type: info.MoveToPosition, Id: 0, Position: info.Position{X: -3215, Y: -1000}},
			},
			{
				{Type: info.MoveToPosition, Id: 1, Position: info.Position{X: -1930, Y: -1000}},
				{Type: info.MoveToPosition, Id: 1, Position: info.Position{X: -1930, Y: 1000}},
			},
			{
				{Type: info.MoveToPosition, Id: 2, Position: info.Position{X: -645, Y: -1000}},
				{Type: info.MoveToPosition, Id: 2, Position: info.Position{X: -645, Y: 1000}},
			}}
	}
	nr_robots := len(instSet)
	return &ObstacleAvoidanceTest{
		team:           team,
		team_size:      team_size,
		max_time:       600 * time.Second,
		at_states:      make([]int, nr_robots),
		instructionSet: instSet,
	}
}

func (m *ObstacleAvoidanceTest) Run() []*info.Instruction {
	var instructions []*info.Instruction

	// This assumes that the robot ids range from 0 to m.team_size
	for id, at_state := range m.at_states {
		instructions = append(instructions, &m.instructionSet[id][at_state])
	}
	return instructions
}

func (m *ObstacleAvoidanceTest) Archived(gs *info.GameState) int {
	// if len(gs.Yellow_team) != 4 {
	// 	fmt.Println("Yellow team has", len(gs.Yellow_team), "4 required")
	// 	return ERROR
	// } else if len(gs.Blue_team) != 3 {
	// 	fmt.Println("Blue team has", len(gs.Blue_team), "3 required")
	// 	return ERROR
	// }

	// This assumes that the robot ids range from 0 to m.team_size
	for id, at_state := range m.at_states {
		target := m.instructionSet[id][at_state].Position
		robot := gs.GetRobot(info.ID(id), m.team)
		if atPosition(robot, target) {
			nr_states := len(m.instructionSet[id])
			m.at_states[id] = (at_state + 1) % nr_states // Cycle through the states
		}
	}
	return RUNNING
}

func atPosition(robot *info.Robot, position info.Position) bool {

	robot_pos := robot.GetPosition()
	diff0 := position.Sub(&robot_pos)
	if diff0.Norm() < 100 {
		return true
	}
	return false
}

// --------------------------MoveToTest----------------------------------------
type MoveToTest struct {
	team     info.Team
	at_state int
	start    time.Time
	max_time time.Duration
}

func NewMoveToTest(team info.Team) *MoveToTest {
	return &MoveToTest{
		team:     team,
		max_time: 20 * time.Second,
		at_state: -1,
	}
}

func (m *MoveToTest) Run() []*info.Instruction {
	if m.at_state == -1 {
		m.start = time.Now()
		m.at_state = 0
	}
	if m.at_state == 0 {
		return []*info.Instruction{
			{Type: info.MoveToPosition, Id: 0, Position: info.Position{X: 1000, Y: 1000}},
			{Type: info.MoveToPosition, Id: 1, Position: info.Position{X: -1000, Y: -1000}},
		}
	} else if m.at_state == 1 {
		return []*info.Instruction{
			{Type: info.MoveToPosition, Id: 0, Position: info.Position{X: 100, Y: 100}},
			{Type: info.MoveToPosition, Id: 1, Position: info.Position{X: -100, Y: -100}},
		}
	} else {
		return []*info.Instruction{
			{Type: info.MoveToPosition, Id: 0, Position: info.Position{X: -1000, Y: -1000}},
			{Type: info.MoveToPosition, Id: 1, Position: info.Position{X: 1000, Y: 1000}},
		}
	}
}

func (m *MoveToTest) Archived(gs *info.GameState) int {
	robot0_pos := gs.GetRobot(info.ID(0), m.team).GetPosition()
	robot1_pos := gs.GetRobot(info.ID(1), m.team).GetPosition()

	if m.at_state == 0 {
		target0 := info.Position{X: 1000, Y: 1000}
		target1 := info.Position{X: -1000, Y: -1000}
		diff0 := target0.Sub(&robot0_pos)
		diff1 := target1.Sub(&robot1_pos)
		if diff0.Norm() < 100 && diff1.Norm() < 100 {
			m.at_state = 1
		}
	} else if m.at_state == 1 {
		target0 := info.Position{X: 100, Y: 100}
		target1 := info.Position{X: -100, Y: -100}
		diff0 := target0.Sub(&robot0_pos)
		diff1 := target1.Sub(&robot1_pos)
		if diff0.Norm() < 100 && diff1.Norm() < 100 {
			// fmt.Println("norms are", diff0.Norm(), diff1.Norm())
			m.at_state = 2
		}
	} else if m.at_state == 2 {
		target0 := info.Position{X: -1000, Y: -1000}
		target1 := info.Position{X: 1000, Y: 1000}
		diff0 := target0.Sub(&robot0_pos)
		diff1 := target1.Sub(&robot1_pos)
		if diff0.Norm() < 100 && diff1.Norm() < 100 {
			m.at_state = 3
		}
	}
	if m.at_state == 3 {
		return COMPLETE
	}
	if m.at_state >= 0 {
		fmt.Println("Time expired", time.Since(m.start))
		if time.Since(m.start) > m.max_time {
			return TIME_EXPIRED
		}
	}

	return RUNNING
}
