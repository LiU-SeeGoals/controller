package ai

import (
	"fmt"
	"time"

	"github.com/LiU-SeeGoals/controller/internal/state"
)

const (
	RUNNING int = iota
	COMPLETE
	TIME_EXPIRED
	ERROR
)

type ScenarioSlowBrain struct {
	team              state.Team
	incomingGameState <-chan state.GameState
	outgoingPlan      chan<- state.GamePlan
	scale             float32
	run_scenario      int // -1 for all
}

func NewScenarioSlowBrain(scale float32, run_scenario int) *ScenarioSlowBrain {
	return &ScenarioSlowBrain{scale: scale, run_scenario: run_scenario}
}

func (sb *ScenarioSlowBrain) Init(incoming <-chan state.GameState, outgoing chan<- state.GamePlan, team state.Team) {
	sb.incomingGameState = incoming
	sb.outgoingPlan = outgoing
	sb.team = team

	go sb.Run()
}

type ScenarioTest interface {
	Run() []*state.Instruction
	Archived(*state.GameState) int
}

func (sb ScenarioSlowBrain) Run() {
	var gameState state.GameState
	gameState.SetValid(false)

	scenarios := []ScenarioTest{}
	scenarios = append(scenarios, NewMoveToTest(sb.team))
	scenarios = append(scenarios, NewObstacleAvoidanceTest(sb.team))
	scenarios = append(scenarios, NewRealTest(sb.team))
	scenario_index := 0
	if sb.run_scenario >= 0 {
		scenario_index = sb.run_scenario
	}

	for {
		gameState = <-sb.incomingGameState

		if !gameState.IsValid() {
			fmt.Println("ScenarioSlowBrain: Invalid game state")
			time.Sleep(10 * time.Millisecond)
			continue
		}

		scenario := scenarios[scenario_index]
		game_state := scenario.Archived(&gameState)
		if game_state == COMPLETE {
			fmt.Println("Scenario", scenario_index, "completed")
		} else if game_state == TIME_EXPIRED {
			fmt.Println("Scenario", scenario_index, "time expired")
		}
		if game_state != RUNNING {
			scenario_index++
			if scenario_index >= len(scenarios) || sb.run_scenario >= 0 {
				panic("ScenarioSlowBrain: No more scenarios") // TODO: Handle this better
			}
			scenario = scenarios[scenario_index]
		}

		inst := scenario.Run()

		plan := state.GamePlan{}
		plan.Instructions = inst

		plan.Team = sb.team

		plan.Valid = true

		sb.outgoingPlan <- plan

	}

}

// --------------------------Real Test--------------------------------
type RealTest struct {
	team           state.Team
	at_state       int
	start          time.Time
	max_time       time.Duration
	instructionSet []*state.Instruction
}

func NewRealTest(team state.Team) *RealTest {

	instSet := []*state.Instruction{
		{Type: state.MoveToPosition, Id: 2, Position: state.Position{X: -2000, Y: 0}},
		{Type: state.MoveToPosition, Id: 2, Position: state.Position{X: 2000, Y: 0}},
	}

	return &RealTest{
		team:           team,
		max_time:       600 * time.Second,
		at_state:       0,
		instructionSet: instSet,
	}
}

func (m *RealTest) Run() []*state.Instruction {

	return []*state.Instruction{m.instructionSet[m.at_state]}
}

func (m *RealTest) Archived(gs *state.GameState) int {
	// if len(gs.Yellow_team) != 4 {
	// 	fmt.Println("Yellow team has", len(gs.Yellow_team), "4 required")
	// 	return ERROR
	// } else if len(gs.Blue_team) != 3 {
	// 	fmt.Println("Blue team has", len(gs.Blue_team), "3 required")
	// 	return ERROR
	// }

	// This assumes that the robot ids range from 0 to m.team_size
	target := m.instructionSet[m.at_state].Position
	robot := gs.GetRobot(state.ID(2), m.team)

	if atPosition(robot, target) {
		m.at_state = (m.at_state + 1) % 2 // Cycle through the states
	}

	return RUNNING
}

// --------------------------Obstacle avoidance--------------------------------
type ObstacleAvoidanceTest struct {
	team           state.Team
	team_size      int
	at_states      []int
	start          time.Time
	max_time       time.Duration
	instructionSet [][]state.Instruction
}

func NewObstacleAvoidanceTest(team state.Team) *ObstacleAvoidanceTest {
	team_size := 4
	var instSet [][]state.Instruction
	if team == state.Yellow {
		instSet = [][]state.Instruction{
			{
				{Type: state.MoveToPosition, Id: 0, Position: state.Position{X: -4000, Y: 0}},
				{Type: state.MoveToPosition, Id: 0, Position: state.Position{X: 4000, Y: 0}},
			},
			{
				{Type: state.MoveToPosition, Id: 1, Position: state.Position{X: 640, Y: -1000}},
				{Type: state.MoveToPosition, Id: 1, Position: state.Position{X: 640, Y: 1000}},
			},
			{ // --------------------------Obstacle avoidance--------------------------------
				{Type: state.MoveToPosition, Id: 2, Position: state.Position{X: 1925, Y: -1000}},
				{Type: state.MoveToPosition, Id: 2, Position: state.Position{X: 1925, Y: 1000}},
			},
			{

				{Type: state.MoveToPosition, Id: 3, Position: state.Position{X: 3210, Y: -1000}},
				{Type: state.MoveToPosition, Id: 3, Position: state.Position{X: 3210, Y: 1000}},
			},
			{

				{Type: state.MoveToPosition, Id: 3, Position: state.Position{X: 3210, Y: -1000}},
				{Type: state.MoveToPosition, Id: 3, Position: state.Position{X: 3210, Y: 1000}},
			},
		}
	} else {
		instSet = [][]state.Instruction{
			{
				{Type: state.MoveToPosition, Id: 0, Position: state.Position{X: -3215, Y: 1000}},
				{Type: state.MoveToPosition, Id: 0, Position: state.Position{X: -3215, Y: -1000}},
			},
			{
				{Type: state.MoveToPosition, Id: 1, Position: state.Position{X: -1930, Y: -1000}},
				{Type: state.MoveToPosition, Id: 1, Position: state.Position{X: -1930, Y: 1000}},
			},
			{
				{Type: state.MoveToPosition, Id: 2, Position: state.Position{X: -645, Y: -1000}},
				{Type: state.MoveToPosition, Id: 2, Position: state.Position{X: -645, Y: 1000}},
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

func (m *ObstacleAvoidanceTest) Run() []*state.Instruction {
	var instructions []*state.Instruction

	// This assumes that the robot ids range from 0 to m.team_size
	for id, at_state := range m.at_states {
		instructions = append(instructions, &m.instructionSet[id][at_state])
	}
	return instructions
}

func (m *ObstacleAvoidanceTest) Archived(gs *state.GameState) int {
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
		robot := gs.GetRobot(state.ID(id), m.team)
		if atPosition(robot, target) {
			nr_states := len(m.instructionSet[id])
			m.at_states[id] = (at_state + 1) % nr_states // Cycle through the states
		}
	}
	return RUNNING
}

func atPosition(robot *state.Robot, position state.Position) bool {

	robot_pos := robot.GetPosition()
	diff0 := position.Sub(&robot_pos)
	if diff0.Norm() < 100 {
		return true
	}
	return false
}

// --------------------------MoveToTest----------------------------------------
type MoveToTest struct {
	team     state.Team
	at_state int
	start    time.Time
	max_time time.Duration
}

func NewMoveToTest(team state.Team) *MoveToTest {
	return &MoveToTest{
		team:     team,
		max_time: 20 * time.Second,
		at_state: -1,
	}
}

func (m *MoveToTest) Run() []*state.Instruction {
	if m.at_state == -1 {
		m.start = time.Now()
		m.at_state = 0
	}
	if m.at_state == 0 {
		return []*state.Instruction{
			{Type: state.MoveToPosition, Id: 0, Position: state.Position{X: 1000, Y: 1000}},
			{Type: state.MoveToPosition, Id: 1, Position: state.Position{X: -1000, Y: -1000}},
		}
	} else if m.at_state == 1 {
		return []*state.Instruction{
			{Type: state.MoveToPosition, Id: 0, Position: state.Position{X: 100, Y: 100}},
			{Type: state.MoveToPosition, Id: 1, Position: state.Position{X: -100, Y: -100}},
		}
	} else {
		return []*state.Instruction{
			{Type: state.MoveToPosition, Id: 0, Position: state.Position{X: -1000, Y: -1000}},
			{Type: state.MoveToPosition, Id: 1, Position: state.Position{X: 1000, Y: 1000}},
		}
	}
}

func (m *MoveToTest) Archived(gs *state.GameState) int {
	robot0_pos := gs.GetRobot(state.ID(0), m.team).GetPosition()
	robot1_pos := gs.GetRobot(state.ID(1), m.team).GetPosition()

	if m.at_state == 0 {
		target0 := state.Position{X: 1000, Y: 1000}
		target1 := state.Position{X: -1000, Y: -1000}
		diff0 := target0.Sub(&robot0_pos)
		diff1 := target1.Sub(&robot1_pos)
		if diff0.Norm() < 100 && diff1.Norm() < 100 {
			m.at_state = 1
		}
	} else if m.at_state == 1 {
		target0 := state.Position{X: 100, Y: 100}
		target1 := state.Position{X: -100, Y: -100}
		diff0 := target0.Sub(&robot0_pos)
		diff1 := target1.Sub(&robot1_pos)
		if diff0.Norm() < 100 && diff1.Norm() < 100 {
			// fmt.Println("norms are", diff0.Norm(), diff1.Norm())
			m.at_state = 2
		}
	} else if m.at_state == 2 {
		target0 := state.Position{X: -1000, Y: -1000}
		target1 := state.Position{X: 1000, Y: 1000}
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
