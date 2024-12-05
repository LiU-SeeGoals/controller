package state

import (
	"encoding/json"
	"fmt"
)

const TEAM_SIZE ID = 11

type GameState struct {
	Valid       bool
	Blue_team   *RobotTeam
	Yellow_team *RobotTeam

	Ball *Ball

	Field *Field

	MessageReceived int64
}

type GameStateDTO struct {
	RobotPositions [2 * TEAM_SIZE]RobotDTO
	BallPosition   BallDTO
}

func (gs GameState) ToDTO() *GameStateDTO {
	gameStateDTO := GameStateDTO{}
	var i ID
	for i = 0; i < TEAM_SIZE; i++ {
		gameStateDTO.RobotPositions[i] = gs.GetRobot(i, Blue).ToDTO()
	}
	for i = 0; i < TEAM_SIZE; i++ {
		gameStateDTO.RobotPositions[TEAM_SIZE+i] = gs.GetRobot(i, Yellow).ToDTO()
	}
	gameStateDTO.BallPosition = gs.Ball.ToDTO()
	return &gameStateDTO
}

func (gs GameState) ToJson() []byte {
	gameStateJson, err := json.Marshal(gs.ToDTO())
	if err != nil {
		fmt.Println("The gamestate packet could not be marshalled to JSON.")
	}
	return gameStateJson
}

func (gs *GameState) SetValid(valid bool) {
	gs.Valid = valid
}

func (gs *GameState) IsValid() bool {
	return gs.Valid
}

func (gs *GameState) SetYellowRobot(robotId uint32, x, y, angle float32, time int64) {
	gs.Yellow_team[robotId].SetPositionTime(x, y, angle, time)

}

func (gs *GameState) SetBlueRobot(robotId uint32, x, y, angle float32, time int64) {
	gs.Blue_team[robotId].SetPositionTime(x, y, angle, time)
}

// Updates position of robots and balls to their actual position
func (gs *GameState) SetBall(x, y, z float32, time int64) {
	gs.Ball.SetPositionTime(x, y, z, time)
}

func (gs *GameState) SetMessageReceivedTime(time int64) {
	gs.MessageReceived = time

}

func (gs *GameState) SetField(length, width, goalWidth, goalDepth, boundaryWidth, penaltyWidth, penaltyDepth int32) {
	gs.Field.FieldLength = length
	gs.Field.FieldWidth = width
	gs.Field.GoalWidth = goalWidth
	gs.Field.GoalDepth = goalDepth
	gs.Field.BoundaryWidth = boundaryWidth
	gs.Field.PenaltyAreaWidth = penaltyWidth
	gs.Field.PenaltyAreaDepth = penaltyDepth
}
func (gs *GameState) AddFieldLine(name string, x1, y1, x2, y2, thickness float32, lineType int) {
	gs.Field.SetLine(name, x1, y1, x2, y2, thickness, FieldShape(lineType))
}

func (gs *GameState) AddFieldArc(name string, centerX, centerY, radius, angle1, angle2, thickness float32, shape int) {
	gs.Field.SetArc(name, centerX, centerY, radius, angle1, angle2, thickness, FieldShape(shape))
}
func (gs *GameState) GetMessageReceivedTime() int64 {
	return gs.MessageReceived
}

func (gs *GameState) GetBall() *Ball {
	return gs.Ball
}

func (gs *GameState) GetTeam(team Team) *RobotTeam {
	if team == Yellow {
		return gs.Yellow_team
	} else {
		return gs.Blue_team
	}
}

func (gs *GameState) GetOtherTeam(team Team) *RobotTeam {
	if team != Yellow {
		return gs.Yellow_team
	} else {
		return gs.Blue_team
	}
}

func (gs *GameState) GetYellowRobots() *RobotTeam {
	return gs.Yellow_team
}

func (gs *GameState) GetBlueRobots() *RobotTeam {
	return gs.Blue_team
}

func (gs *GameState) GetRobot(id ID, team Team) *Robot {
	if team == Blue {
		return gs.Blue_team[id]
	}
	return gs.Yellow_team[id]
}

// Constructor for game state
func NewGameState(capacity int) *GameState {
	gs := GameState{}
	gs.Valid = true

	gs.Ball = NewBall(capacity)
	var i ID
	gs.Blue_team = new(RobotTeam)
	gs.Yellow_team = new(RobotTeam)
	for i = 0; i < TEAM_SIZE; i++ {
		gs.Blue_team[i] = NewRobot(i, Blue, capacity)
		gs.Yellow_team[i] = NewRobot(i, Yellow, capacity)
	}
	gs.Field = NewField()

	return &gs

}

// String representation of game state
func (gs *GameState) String() string {
	gs_str := "{\n blue team: {\n"
	var i ID
	for i = 0; i < TEAM_SIZE; i++ {
		gs_str += "robot: {" + gs.Blue_team[i].String() + " },\n"
	}
	gs_str += "},\n yellow team: {\n"
	for i = 0; i < TEAM_SIZE; i++ {
		gs_str += "robot: {" + gs.Yellow_team[i].String() + " },\n"
	}
	// for _, line := range gs.Field.FieldLines {
	// 	gs_str += fmt.Sprintf("line: {%s}\n", line.String())
	// }
	// for _, arc := range gs.Field.FieldArcs {
	// 	gs_str += fmt.Sprintf("arc: {%s}\n", arc.String())
	// }
	gs_str += "}"
	return gs_str
}
