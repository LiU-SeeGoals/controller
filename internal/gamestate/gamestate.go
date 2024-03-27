package gamestate

import (
	"fmt"
)

const TEAM_SIZE = 6

type GameState struct {
	Blue_team   [TEAM_SIZE]*Robot
	Yellow_team [TEAM_SIZE]*Robot

	// Holds ball data
	ball *Ball
	// Holds field data
	Field Field
}

func (gs *GameState) SetRobot(robotId uint32, x, y, w float64, isBlue bool) {
	if isBlue {
		gs.Blue_team[robotId].SetPosition(x, y, w)
	} else {
		gs.Yellow_team[robotId].SetPosition(x, y, w)
	}
}

func (gs *GameState) SetBall(x, y, z float64) {
	gs.ball.SetPosition(x, y, z)
}

func (gs *GameState) GetBall() *Ball {
	return gs.ball
}

func (gs *GameState) GetTeam(team Team) [TEAM_SIZE]*Robot {
	if team == Yellow {
		return gs.Yellow_team
	} else {
		return gs.Blue_team
	}
}

func (gs *GameState) GetRobot(id int, isBlue bool) *Robot {
	if isBlue {
		return gs.Blue_team[id]
	}
	return gs.Yellow_team[id]
}

func NewGameState() *GameState {
	gs := &GameState{}

	gs.ball = NewBall()

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
