package ai

import (
	"bytes"

	"encoding/json"
	"io"
	"net/http"
	"time"

	. "github.com/LiU-SeeGoals/controller/internal/logger"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type SlowBrainPy struct {
	team              info.Team
	incomingGameState <-chan info.GameState
	outgoingPlan      chan<- info.GamePlan
	ip_address        string
}

func NewSlowBrainPy(ip_address string) *SlowBrainPy {
	return &SlowBrainPy{ip_address: ip_address}
}

func (sb *SlowBrainPy) Init(incoming <-chan info.GameState, outgoing chan<- info.GamePlan, team info.Team) {
	sb.incomingGameState = incoming
	sb.outgoingPlan = outgoing
	sb.team = team

	go sb.Run()
}

type PyResponse struct {
	Instructions []struct {
		Id       int
		Position []float32
	}
}

func (sb SlowBrainPy) Run() {
	var gameState info.GameState
	gameState.SetValid(false)

	for {
		gameState = <-sb.incomingGameState

		if !gameState.IsValid() {
			// fmt.Println("ScenarioSlowBrain: Invalid game state")
			Logger.Warn("ScenarioSlowBrain: Invalid game state")
			time.Sleep(10 * time.Millisecond)
			continue
		}

		resp, err := http.Post(sb.ip_address, "application/json",
			bytes.NewBuffer(gameState.ToJson()))
		if err != nil {
			// fmt.Println("Error in http.Get")
			// fmt.Println(err)
			Logger.Errorf("Error in http.Get: %v", err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			// fmt.Println("Error in io.ReadAll")
			// fmt.Println(err)
			Logger.Errorf("Error in io.ReadAll: %v", err)
			continue
		}

		var pyResponse PyResponse

		err = json.Unmarshal(body, &pyResponse)

		if err != nil {
			// fmt.Println("Error in json.Unmarshal")
			// fmt.Println(err)
			Logger.Errorf("Error in json.Unmarshal: %v", err)
			continue
		}

		// loop over the json["instructions"] response in the body and create a plan

		plan := info.GamePlan{}
		for _, scenario := range pyResponse.Instructions {
			idx := scenario.Id

			robot := gameState.GetRobot(info.ID(idx), sb.team)
			position := info.Position{
				X:     scenario.Position[0],
				Y:     scenario.Position[1],
				Z:     scenario.Position[2],
				Angle: scenario.Position[3],
			}
			plan.Instructions = append(plan.Instructions, &info.Instruction{
				Type:     info.MoveToPosition,
				Id:       robot.GetID(),
				Position: position,
			})
		}

		plan.Team = sb.team

		plan.Valid = true

		sb.outgoingPlan <- plan

	}

}
