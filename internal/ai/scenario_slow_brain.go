package ai

import (
	"fmt"
	"time"

	"github.com/LiU-SeeGoals/controller/internal/state"
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

type MoveToTest struct {
	team     state.Team
	at_state int // 0 for init, 1 for goal
}

func (m *MoveToTest) Run() []*state.Instruction {
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

func (m *MoveToTest) Archived(gs *state.GameState) bool {
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
	return m.at_state == 3
}

func (sb ScenarioSlowBrain) Run() {
	var gameState state.GameState
	gameState.SetValid(false)

	scenarios := []MoveToTest{}
	scenarios = append(scenarios, MoveToTest{team: sb.team, at_state: 0})
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

		scenario := &scenarios[scenario_index]

		if scenario.Archived(&gameState) {
			scenario_index++
			if scenario_index >= len(scenarios) || sb.run_scenario >= 0 {
				panic("ScenarioSlowBrain: No more scenarios") // TODO: Handle this better
			}
			scenario = &scenarios[scenario_index]
		}

		inst := scenario.Run()

		plan := state.GamePlan{}
		plan.Instructions = inst

		plan.Team = sb.team

		plan.Valid = true

		sb.outgoingPlan <- plan

	}

}
