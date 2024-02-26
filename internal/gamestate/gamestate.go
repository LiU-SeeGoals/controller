package gamestate

import (
	"fmt"
)

const TEAM_SIZE = 11

type GameState struct {
	blue_team   [TEAM_SIZE]*Robot
	yellow_team [TEAM_SIZE]*Robot

	// Holds ball data
	ball *Ball
	// Holds field data
	Field Field
}

func (gs *GameState) SetRobot(robotId uint32, x, y, w float64, team Team) {
	if team == Blue {
		gs.blue_team[robotId].SetPosition(x, y)
		gs.blue_team[robotId].SetAngle(w)
	} else {
		gs.yellow_team[robotId].SetPosition(x, y)
		gs.blue_team[robotId].SetAngle(w)
	}
}

func (gs *GameState) SetBall(x, y, z float64) {
	gs.ball.SetPosition(x, y, z)
}

func (gs *GameState) GetBall() *Ball {
	return gs.ball
}

func (gs *GameState) GetRobot(id int, team Team) *Robot {
	if team == Blue {
		return gs.blue_team[id]
	}
	return gs.yellow_team[id]
}

func NewGameState() *GameState {
	gs := &GameState{}

	gs.ball = NewBall()

	for i := 0; i < TEAM_SIZE; i++ {
		gs.blue_team[i] = NewRobot(i, Blue)
		gs.yellow_team[i] = NewRobot(i, Yellow)
	}

	return gs
}

// String representation of game state
func (gs *GameState) String() string {
	gs_str := "{\n blue team: {\n"
	for i := 0; i < TEAM_SIZE; i++ {
		gs_str += "robot: {" + gs.blue_team[i].String() + " },\n"
	}
	gs_str += "},\n yellow team: {\n"
	for i := 0; i < TEAM_SIZE; i++ {
		gs_str += "robot: {" + gs.yellow_team[i].String() + " },\n"
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
