package gamestate

import (
	"fmt"

	"github.com/LiU-SeeGoals/controller/internal/config"
	"github.com/LiU-SeeGoals/controller/internal/parsed_vision"
)

type GameState struct {
	frameNumber int
	blue_team   map[int]*Robot
	yellow_team map[int]*Robot

	// Holds ball data
	ball *Ball
	// Holds field data
	Field Field
}

func (gs *GameState) SetRobot(robotId int, x, y, w float64, isBlue bool) {
	if isBlue {
		gs.blue_team[robotId].SetPosition(x, y, w)
	} else {
		gs.yellow_team[robotId].SetPosition(x, y, w)
	}
}

func (gs *GameState) SetBall(x, y, z float64) {
	gs.ball.SetPosition(x, y, z)
}

func (gs *GameState) GetBall() *Ball {
	return gs.ball
}

func (gs *GameState) GetFrameNumber() int {
	return gs.frameNumber
}

func (gs *GameState) SetFrameNumber(frameNumber int) {
	gs.frameNumber = frameNumber
}

func (gs *GameState) AddRobot(id int, isBlue bool) {
	if isBlue {
		gs.blue_team[id] = NewRobot(id, Blue)
		return
	}
	gs.yellow_team[id] = NewRobot(id, Yellow)
}

func (gs *GameState) GetRobot(id int, isBlue bool) *Robot {
	if isBlue {
		return gs.blue_team[id]
	}
	return gs.yellow_team[id]
}

func (gs *GameState) Clear() {
	gs.blue_team = make(map[int]*Robot)
	gs.yellow_team = make(map[int]*Robot)
	gs.ball = NewBall()
	gs.Field = Field{}
}

func (gs *GameState) GetParsedGameState() *parsed_vision.ParsedFrame {
	var robots []*parsed_vision.Robot
	var parsedFrame = parsed_vision.ParsedFrame{}

	parsedFrame.Ball = gs.ball.GetParsedBall()

	isBlue := config.GetIsBlueTeam() // Subject to change (idk how we get info on which team color we play as.)
	selectedTeam := gs.yellow_team
	if isBlue {
		selectedTeam = gs.blue_team
	}

	for _, robot := range selectedTeam {
		robots = append(robots, robot.GetParsedRobot())
	}

	parsedFrame.Robots = robots 

	return &parsedFrame
}

func NewGameState() *GameState {
	gs := &GameState{}

	gs.ball = NewBall()

	gs.blue_team = make(map[int]*Robot)
	gs.yellow_team = make(map[int]*Robot)

	return gs
}

// String representation of game state
func (gs *GameState) String() string {
	gs_str := "{\n blue team: {\n"
	for i, _ := range gs.blue_team{
		gs_str += "robot: {" + gs.blue_team[i].String() + " },\n"
	}
	gs_str += "},\n yellow team: {\n"
	for i, _ := range gs.yellow_team {
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

