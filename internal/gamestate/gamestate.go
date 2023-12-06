package gamestate

import (
	"fmt"

	"github.com/LiU-SeeGoals/controller/internal/action"

	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/proto/ssl_vision"
	"github.com/LiU-SeeGoals/controller/internal/receiver"
	"gonum.org/v1/gonum/mat"
)

const TEAM_SIZE = 6

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
}

// Method used for testing actions,
// a proper test should be implemented
func (gs *GameState) TestActions() {
	id := 0
	act := &action.Move{}
	act.Pos = gs.yellow_team[id].pos
	act.Dest = mat.NewVecDense(3, nil)
	act.Id = id
	act.Dest.SetVec(0, 0)
	act.Dest.SetVec(1, 0)
	act.Dest.SetVec(2, 0)
	//act.Dribble = true

	//act := &action.Kick{}
	//act.Kickspeed = 10

	//act := &action.Dribble{}
	//act.Dribble = true

	var action []action.Action
	action = append(action, act)

	gs.Grsim_client.SendActions(action)
}

// Updates position of robots and balls to their actual position
func (gs *GameState) Update() {
	var packet ssl_vision.SSL_WrapperPacket

	var detect *ssl_vision.SSL_DetectionFrame
	var field *ssl_vision.SSL_GeometryFieldSize

	packet = <-gs.ssl_receiver_channel

	detect = packet.GetDetection()

	geo := packet.GetGeometry()
	if geo != nil {
		field = geo.GetField()
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

	parseFieldData(&gs.field, field)
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
func parseFieldData(f *Field, data *ssl_vision.SSL_GeometryFieldSize) {
	if data == nil {
		return
	}

	// parse field data
	f.FieldLengt = data.GetFieldLength()
	f.FieldWidth = data.GetFieldWidth()
	f.BallRadius = data.GetBallRadius()
	f.BoundaryWidth = data.GetBoundaryWidth()
	f.CenterRadius = data.GetCenterCircleRadius()
	f.GoalDepth = data.GetGoalDepth()
	f.GoalHeight = data.GetGoalHeight()
	f.GoalWidth = data.GetGoalWidth()
	f.GoalToPenalty = data.GetGoalCenterToPenaltyMark()
	f.LineThickness = data.GetLineThickness()
	f.MaxRobotRadius = data.GetMaxRobotRadius()
	f.PenaltyAreaDepth = data.GetPenaltyAreaDepth()
	f.PenaltyAreaWidth = data.GetPenaltyAreaWidth()

	parseFieldLines(f, data.GetFieldLines())
	parseFieldArcs(f, data.GetFieldArcs())
}

// Parse field lines from ssl packet
//
// Field object should be passed from game state object.
func parseFieldLines(f *Field, lines []*ssl_vision.SSL_FieldLineSegment) {
	for _, line := range lines {
		if f.hasLine(line.GetName()) {
			continue
		}
		p1 := line.GetP1()
		p2 := line.GetP2()
		f.addLine(
			line.GetName(),
			p1.GetX(),
			p1.GetY(),
			p2.GetX(),
			p2.GetY(),
			line.GetThickness(),
			convertShapeType(line.GetType()),
		)
	}
}

// Parse arcs from ssl packet
//
// Field object should be passed from game state object.
func parseFieldArcs(f *Field, arcs []*ssl_vision.SSL_FieldCircularArc) {
	for _, arc := range arcs {
		if f.hasArc(arc.GetName()) {
			continue
		}

		center := arc.GetCenter()
		f.addArc(
			arc.GetName(),
			center.GetX(),
			center.GetY(),
			arc.GetRadius(),
			arc.GetA1(),
			arc.GetA2(),
			arc.GetThickness(),
			convertShapeType(arc.GetType()),
		)
	}
}

// Glorified type cast
// Converts ssl vision enum to our own enum
func convertShapeType(typ ssl_vision.SSL_FieldShapeType) FieldShape {
	return FieldShape(typ)
}
