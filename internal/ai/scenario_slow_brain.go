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
}

func NewScenarioSlowBrain(scale float32) ScenarioSlowBrain {
	return ScenarioSlowBrain{scale: scale}
}

func (sb *ScenarioSlowBrain) Init(incoming <-chan state.GameState, outgoing chan<- state.GamePlan, team state.Team) {
	sb.incomingGameState = incoming
	sb.outgoingPlan = outgoing
	sb.team = team

	go sb.Run()
}

func (sb ScenarioSlowBrain) scenarioArchived(gameState *state.GameState, scenario []state.Position) bool {
	for idx, scenario_pos := range scenario {
		robot := gameState.GetRobot(state.ID(idx), sb.team)
		pos := robot.GetPosition()
		diff := scenario_pos.Sub(&pos)
		if diff.Norm() > 100 {
			return false
		}
	}
	return true
}

func (sb ScenarioSlowBrain) Run() {
	var gameState state.GameState
	gameState.SetValid(false)
	scenario_index := 0
	scenarios := [][]state.Position{
		{
			{X: sb.scale * 100, Y: sb.scale * 100},
			{X: sb.scale * 200, Y: sb.scale * 200},
			{X: sb.scale * 300, Y: sb.scale * 300},
		},
		{
			{X: sb.scale * 100, Y: sb.scale * -100},
			{X: sb.scale * 200, Y: sb.scale * -200},
			{X: sb.scale * 300, Y: sb.scale * -300},
		},
		{
			{X: sb.scale * -100, Y: sb.scale * -100},
			{X: sb.scale * -200, Y: sb.scale * -200},
			{X: sb.scale * -300, Y: sb.scale * -300},
		},
		{
			{X: sb.scale * -100, Y: sb.scale * 100},
			{X: sb.scale * -200, Y: sb.scale * 200},
			{X: sb.scale * -300, Y: sb.scale * 300},
		},
	}

	for {
		gameState = <-sb.incomingGameState

		if !gameState.IsValid() {
			fmt.Println("ScenarioSlowBrain: Invalid game state")
			time.Sleep(10 * time.Millisecond)
			continue
		}

		scenario := scenarios[scenario_index]

		if sb.scenarioArchived(&gameState, scenario) {
			fmt.Println("Scenario archived: ", scenario)
			scenario_index = (scenario_index + 1) % len(scenarios)
			scenario = scenarios[scenario_index]
		}

		plan := state.GamePlan{}
		for idx, scenario_pos := range scenario {
			robot := gameState.GetRobot(state.ID(idx), sb.team)
			plan.Instructions = append(plan.Instructions, state.RobotMove{
				Id:       robot.GetID(),
				Position: scenario_pos,
			})
		}

		plan.Team = sb.team

		plan.Valid = true

		sb.outgoingPlan <- plan

	}

}
