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

	actions []action.Action
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
	gameStateDTO.Actions = make([]action.ActionDTO, len(gs.actions))

	for i, act := range gs.actions {
		gameStateDTO.Actions[i] = act.ToDTO()
	}

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

func (gs *GameState) SetBall(x, y, z float64) {
	gs.Ball.SetPosition(x, y, z)
}

//func (gs *GameState) AddAction(action action.Action) {
//	gs.actions = append(gs.actions, action)
//}
//
//func (gs *GameState) sendActions() {
//	gs.Grsim_client.SendActions(gs.actions)
//	gs.actions = nil
//}

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

func NewGameState(sslClientAddress string, sslReceiverAddress string) *GameState {
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

// Parse geoemtry field data
// func parseFieldData(f *Field, data *ssl_vision.SSL_GeometryFieldSize) {
// 	if data == nil {
// 		return
// 	}

// 	// parse field data
// 	f.FieldLengt = data.GetFieldLength()
// 	f.FieldWidth = data.GetFieldWidth()
// 	f.BallRadius = data.GetBallRadius()
// 	f.BoundaryWidth = data.GetBoundaryWidth()
// 	f.CenterRadius = data.GetCenterCircleRadius()
// 	f.GoalDepth = data.GetGoalDepth()
// 	f.GoalHeight = data.GetGoalHeight()
// 	f.GoalWidth = data.GetGoalWidth()
// 	f.GoalToPenalty = data.GetGoalCenterToPenaltyMark()
// 	f.LineThickness = data.GetLineThickness()
// 	f.MaxRobotRadius = data.GetMaxRobotRadius()
// 	f.PenaltyAreaDepth = data.GetPenaltyAreaDepth()
// 	f.PenaltyAreaWidth = data.GetPenaltyAreaWidth()

// 	parseFieldLines(f, data.GetFieldLines())
// 	parseFieldArcs(f, data.GetFieldArcs())
// }

// Parse field lines from ssl packet
//
// Field object should be passed from game state object.
// func parseFieldLines(f *Field, lines []*ssl_vision.SSL_FieldLineSegment) {
// 	for _, line := range lines {
// 		if f.hasLine(line.GetName()) {
// 			continue
// 		}
// 		p1 := line.GetP1()
// 		p2 := line.GetP2()
// 		f.addLine(
// 			line.GetName(),
// 			p1.GetX(),
// 			p1.GetY(),
// 			p2.GetX(),
// 			p2.GetY(),
// 			line.GetThickness(),
// 			convertShapeType(line.GetType()),
// 		)
// 	}
// }

// // Parse arcs from ssl packet
// //
// // Field object should be passed from game state object.
// func parseFieldArcs(f *Field, arcs []*ssl_vision.SSL_FieldCircularArc) {
// 	for _, arc := range arcs {
// 		if f.hasArc(arc.GetName()) {
// 			continue
// 		}

// 		center := arc.GetCenter()
// 		f.addArc(
// 			arc.GetName(),
// 			center.GetX(),
// 			center.GetY(),
// 			arc.GetRadius(),
// 			arc.GetA1(),
// 			arc.GetA2(),
// 			arc.GetThickness(),
// 			convertShapeType(arc.GetType()),
// 		)
// 	}
// }
