package gamestate

import (
	"encoding/json"
	"fmt"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/webserver"

	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/receiver"
	"github.com/LiU-SeeGoals/proto_go/ssl_vision"
)

const TEAM_SIZE = 12

type GameState struct {
	Grsim_client         *client.GrsimClient
	ssl_receiver         *receiver.SSLReceiver
	ssl_receiver_channel chan ssl_vision.SSL_WrapperPacket

	blue_team   [TEAM_SIZE]*Robot
	yellow_team [TEAM_SIZE]*Robot

	// Holds ball data
	ball *Ball
	// Holds field data
	field Field

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

	gameStateDTO.Ball = gs.ball.ToDTO()
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

func (gs *GameState) AddAction(action action.Action) {
	gs.actions = append(gs.actions, action)
}

func (gs *GameState) sendActions() {
	gs.Grsim_client.SendActions(gs.actions)
	gs.actions = nil
}

// Updates position of robots and balls to their actual position
func (gs *GameState) Update() {
	// helloman = helloman + 1
	// fmt.Println(helloman)
	var packet ssl_vision.SSL_WrapperPacket

	var detect *ssl_vision.SSL_DetectionFrame
	//var field *ssl_vision.SSL_GeometryFieldSize

	packet = <-gs.ssl_receiver_channel

	detect = packet.GetDetection()

	geo := packet.GetGeometry()
	if geo != nil {
		//gs.field = geo.GetField()
	}

	for _, robot := range detect.GetRobotsBlue() {
		x := float64(robot.GetX())
		y := float64(robot.GetY())
		w := float64(*robot.Orientation)

		gs.blue_team[robot.GetRobotId()].SetPosition(x, y, w)
	}

	for _, robot := range detect.GetRobotsYellow() {
		x := float64(robot.GetX())
		y := float64(robot.GetY())
		w := float64(*robot.Orientation)

		gs.yellow_team[robot.GetRobotId()].SetPosition(x, y, w)

	}

	for _, ball := range detect.GetBalls() {
		x := float64(ball.GetX())
		y := float64(ball.GetY())
		w := float64(ball.GetZ())

		gs.ball.SetPosition(x, y, w)
	}

	//parseFieldData(&gs.field, field)
	gs.broadcastGameState()
	gs.sendActions()
}

func (gs *GameState) broadcastGameState() {
	webserver.BroadcastGameState(gs.ToJson())
	// list of incoming actions
	incomming := webserver.GetIncoming()

	// If we got new actions --> then handle them
	if len(incomming) > 0 {
		gs.handleIncoming(incomming)
	}
}

func (gs *GameState) GetBall() *Ball {
	return gs.ball
}

func (gs *GameState) handleIncoming(incomming []action.ActionDTO) {
	fmt.Println("Received a new action (gamestate)")
	// for _, act := range incomming {
	// 	switch act.Action {
	// 	case robot_action.ActionType_MOVE_ACTION:
	// 		pos := mat.NewVecDense(3, []float64{float64(act.PosX), float64(act.PosY), float64(act.PosW)})
	// 		dest := mat.NewVecDense(3, []float64{float64(act.DestX), float64(act.DestY), float64(act.DestW)})
	// 		gs.AddAction(&action.Move{act.Id, pos, dest, act.Dribble})
	// 	case robot_action.ActionType_INIT_ACTION:
	// 		gs.AddAction(&action.Init{act.Id})
	// 	case robot_action.ActionType_ROTATE_ACTION:
	// 		gs.AddAction(&action.Rotate{act.Id, int(act.PosW)})
	// 	case robot_action.ActionType_KICK_ACTION:
	// 		standardKickSpeed := 1
	// 		gs.AddAction(&action.Kick{act.Id, standardKickSpeed})
	// 	case robot_action.ActionType_MOVE_TO_ACTION:
	// 		dest := mat.NewVecDense(3, []float64{float64(act.DestX), float64(act.DestY)})
	// 		gs.AddAction(&action.SetNavigationDirection{act.Id, dest})
	// 	case robot_action.ActionType_STOP_ACTION:
	// 		gs.AddAction(&action.Stop{act.Id})
	// 	}
	// }

}

func (gs *GameState) GetRobot(id int, team Team) *Robot {
	if team == Blue {
		return gs.blue_team[id]
	}
	return gs.yellow_team[id]
}

// Start a SSL Vision receiver, returns a channel from
// which SSL wrapper packets can be obtained.
func (gs *GameState) setupSSLVisionReceiver(addr string) {
	gs.ssl_receiver = receiver.NewSSLReceiver(addr)
	gs.ssl_receiver.Connect()

	gs.ssl_receiver_channel = make(chan ssl_vision.SSL_WrapperPacket)
	go gs.ssl_receiver.Receive(gs.ssl_receiver_channel)
}

func NewGameState(sslClientAddress string, sslReceiverAddress string) *GameState {
	gs := &GameState{}

	gs.Grsim_client = client.NewGrsimClient(sslClientAddress)
	gs.Grsim_client.Init()

	gs.setupSSLVisionReceiver(sslReceiverAddress)

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
	for _, line := range gs.field.FieldLines {
		gs_str += fmt.Sprintf("line: {%s}\n", line.String())
	}
	for _, arc := range gs.field.FieldArcs {
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

// Glorified type cast
// Converts ssl vision enum to our own enum
func convertShapeType(typ ssl_vision.SSL_FieldShapeType) FieldShape {
	return FieldShape(typ)
}
