package gamestate

import (
	"encoding/json"
	"fmt"
)

const TEAM_SIZE = 6

type GameState struct {
	Blue_team   [TEAM_SIZE]*Robot
	Yellow_team [TEAM_SIZE]*Robot

	// Holds ball data
	Ball *Ball

	MessageReceived int64
	LagTime         int64
}

func (gs *GameState) copy(clone *GameState) {
	gs.Ball.copy(clone.Ball)
	for i := 0; i < TEAM_SIZE; i++ {
		gs.Blue_team[i].copy(clone.Blue_team[i])
		gs.Yellow_team[i].copy(clone.Yellow_team[i])
	}
	clone.MessageReceived = gs.MessageReceived
	clone.LagTime = gs.LagTime
}

type GameStateDTO struct {
	RobotPositions [2 * TEAM_SIZE]RobotDTO
	BallPosition   BallDTO
}

func (gs *GameState) ToDTO() *GameStateDTO {
	gameStateDTO := GameStateDTO{}
	for i := 0; i < TEAM_SIZE; i++ {
		gameStateDTO.RobotPositions[i] = gs.GetRobot(i, Blue).ToDTO()
	}
	for i := 0; i < TEAM_SIZE; i++ {
		gameStateDTO.RobotPositions[TEAM_SIZE+i] = gs.GetRobot(i, Yellow).ToDTO()
	}
	gameStateDTO.BallPosition = gs.Ball.ToDTO()
	return &gameStateDTO
}

func (gs *GameState) ToJson() []byte {
	gameStateJson, err := json.Marshal(gs.ToDTO())
	if err != nil {
		fmt.Println("The gamestate packet could not be marshalled to JSON.")
	}
	return gameStateJson
}

func (gs *GameState) SetYellowRobot(robotId uint32, x, y, w float64, time int64) {
	gs.Yellow_team[robotId].SetPositionTime(x, y, w, time)
	gs.Yellow_team[robotId].UpdateVelocity()
}

func (gs *GameState) SetBlueRobot(robotId uint32, x, y, w float64, time int64) {
	gs.Blue_team[robotId].SetPositionTime(x, y, w, time)
	gs.Blue_team[robotId].UpdateVelocity()
}

func (gs *GameState) UpdateYellowRobot(robotId uint32, x, y, w float64, time int64) {
	gs.Yellow_team[robotId].UpdatePositionTime(x, y, w, time)
}

func (gs *GameState) UpdateBlueRobot(robotId uint32, x, y, w float64, time int64) {
	gs.Blue_team[robotId].UpdatePositionTime(x, y, w, time)
}

// Updates position of robots and balls to their actual position
func (gs *GameState) SetBall(x, y, z float64, time int64) {
	gs.Ball.SetPositionTime(x, y, z, time)
	gs.Ball.UpdateVelocity()
}

func (gs *GameState) SetMessageReceivedTime(time int64) {
	gs.MessageReceived = time

}
func (gs *GameState) GetMessageReceivedTime() int64 {
	return gs.MessageReceived
}

func (gs *GameState) SetLagTime(lagTime int64) {
	gs.LagTime = lagTime
}

func (gs *GameState) GetLagTime() int64 {
	return gs.LagTime
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

func (gs *GameState) GetYellowRobots() [TEAM_SIZE]*Robot {
	return gs.Yellow_team
}

func (gs *GameState) GetBlueRobots() [TEAM_SIZE]*Robot {
	return gs.Blue_team
}

func (gs *GameState) GetRobot(id int, team Team) *Robot {
	if team == Blue {
		return gs.Blue_team[id]
	}
	return gs.Yellow_team[id]
}

// Constructor for game state
func NewGameState(capasity int) *GameState {
	gs := &GameState{}

	gs.Ball = NewBall(capasity)
	for i := 0; i < TEAM_SIZE; i++ {
		gs.Blue_team[i] = NewRobot(i, Blue, capasity)
		gs.Yellow_team[i] = NewRobot(i, Yellow, capasity)
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
