package gamestate

import (
	"encoding/json"
	"fmt"

	"github.com/LiU-SeeGoals/controller/internal/action"
)

const TEAM_SIZE = 12

type GameState struct {
	Blue_team   [TEAM_SIZE]*Robot
	Yellow_team [TEAM_SIZE]*Robot

	// Holds ball data
	Ball *Ball
	// Holds field data
	Field *Field
}

type GameStateDTO struct {
	BlueTeam   [TEAM_SIZE]RobotDTO
	YellowTeam [TEAM_SIZE]RobotDTO
	Ball       BallDTO
	Actions    []action.ActionDTO
}

func (gs *GameState) ToDTO() *GameStateDTO {

	gameStateDTO := GameStateDTO{}

	for i := 0; i < TEAM_SIZE; i++ {
		gameStateDTO.BlueTeam[i] = gs.GetRobot(i, Blue).ToDTO()
		gameStateDTO.YellowTeam[i] = gs.GetRobot(i, Yellow).ToDTO()
	}

	gameStateDTO.Ball = gs.Ball.ToDTO()
	// gameStateDTO.Actions = make([]action.ActionDTO, len(gs.ManualActions))

	// for i, act := range gs.ManualActions {
	// 	gameStateDTO.Actions[i] = act.ToDTO()
	// }

	return &gameStateDTO
}

func (gs *GameState) ToJson() []byte {
	gameStateJson, err := json.Marshal(gs.ToDTO())
	if err != nil {
		fmt.Println("The gamestate packet could not be marshalled to JSON.")
	}
	return gameStateJson
}

func (gs *GameState) SetRobot(robotId uint32, x, y, w float64, isBlue bool) {
	if isBlue {
		gs.Blue_team[robotId].SetPosition(x, y, w)
	} else {
		gs.Yellow_team[robotId].SetPosition(x, y, w)
	}
}

// Updates position of robots and balls to their actual position
func (gs *GameState) SetBall(x, y, z float64) {
	gs.Ball.SetPosition(x, y, z)
}

func (gs *GameState) GetBall() *Ball {
	return gs.Ball
}

func (gs *GameState) GetTeam(team Team) [TEAM_SIZE]*Robot {
	if team == Yellow {
		return gs.Yellow_team
	} else {
		return gs.Blue_team
	}
}

func (gs *GameState) GetRobot(id int, team Team) *Robot {
	if team == Blue {
		return gs.Blue_team[id]
	}
	return gs.Yellow_team[id]
}

func NewGameState() *GameState {
	gs := &GameState{}

	gs.Ball = NewBall()
	for i := 0; i < TEAM_SIZE; i++ {
		gs.Blue_team[i] = NewRobot(i, Blue)
		gs.Yellow_team[i] = NewRobot(i, Yellow)
	}

	return gs

}

// String representation of game state
func (gs *GameState) String() string {
	gs_str := "{\n blue team: {\n"
	for i := 0; i < TEAM_SIZE; i++ {
		gs_str += "robot: {" + gs.Blue_team[i].String() + " },\n"
	}
	gs_str += "},\n yellow team: {\n"
	for i := 0; i < TEAM_SIZE; i++ {
		gs_str += "robot: {" + gs.Yellow_team[i].String() + " },\n"
	}
	for _, line := range gs.Field.FieldLines {
		gs_str += fmt.Sprintf("line: {%s}\n", line.String())
	}
	for _, arc := range gs.Field.FieldArcs {
		gs_str += fmt.Sprintf("arc: {%s}\n", arc.String())
	}
	gs_str += "}"
	return gs_str
}
