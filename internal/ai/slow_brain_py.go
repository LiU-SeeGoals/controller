package ai

import (
	"bytes"
	"fmt"

	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/LiU-SeeGoals/controller/internal/state"
)

type SlowBrainPy struct {
	team              state.Team
	incomingGameState <-chan state.GameState
	outgoingPlan      chan<- state.GamePlan
	ip_address        string
}

func NewSlowBrainPy(ip_address string) SlowBrainPy {
	return SlowBrainPy{ip_address: ip_address}
}

func (sb *SlowBrainPy) Init(incoming <-chan state.GameState, outgoing chan<- state.GamePlan, team state.Team) {
	sb.incomingGameState = incoming
	sb.outgoingPlan = outgoing
	sb.team = team

	go sb.Run()
}

// resp = {
// 	"instructions": [
// 		{"id": 0,
// 		 "position": [0, 0, 0, 0],
// 		},
// 		{"id": 1,
// 		 "position": [0, 0, 0, 0],
// 		},
// 	]

// }

type PyResponse struct {
	Instructions []struct {
		Id       int
		Position []float32
	}
}

func (sb SlowBrainPy) Run() {
	var gameState state.GameState

	for {
		gameState = <-sb.incomingGameState

		if !gameState.IsValid() {
			fmt.Println("ScenarioSlowBrain: Invalid game state")
			time.Sleep(10 * time.Millisecond)
			continue
		}

		resp, err := http.Post(sb.ip_address, "application/json",
			bytes.NewBuffer(gameState.ToJson()))
		if err != nil {
			fmt.Println("Error in http.Get")
			fmt.Println(err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error in io.ReadAll")
			fmt.Println(err)
			continue
		}

		var pyResponse PyResponse

		err = json.Unmarshal(body, &pyResponse)

		if err != nil {
			fmt.Println("Error in json.Unmarshal")
			fmt.Println(err)
			continue
		}

		// loop over the json["instructions"] response in the body and create a plan

		plan := state.GamePlan{}
		for _, scenario := range pyResponse.Instructions {
			idx := scenario.Id

			robot := gameState.GetRobot(state.ID(idx), sb.team)
			position := state.Position{
				X:     scenario.Position[0],
				Y:     scenario.Position[1],
				Z:     scenario.Position[2],
				Angle: scenario.Position[3],
			}
			plan.Instructions = append(plan.Instructions, state.RobotMove{
				Id:       robot.GetID(),
				Position: position,
			})
		}

		plan.Team = sb.team

		plan.Valid = true

		sb.outgoingPlan <- plan

	}

}
