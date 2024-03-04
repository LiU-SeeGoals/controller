package world_predictor

import (
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"github.com/LiU-SeeGoals/controller/internal/receiver"
	"github.com/LiU-SeeGoals/proto-messages/ssl_vision"
)

type WorldPredictor struct {
	ssl_receiver         *receiver.SSLReceiver
	ssl_receiver_channel chan ssl_vision.SSL_WrapperPacket
	gamestate            *gamestate.GameState
}

// Updates position of robots and balls to their actual position
func (wp *WorldPredictor) Update() {
	var packet ssl_vision.SSL_WrapperPacket

	var detect *ssl_vision.SSL_DetectionFrame
	var field *ssl_vision.SSL_GeometryFieldSize

	packet = <-wp.ssl_receiver_channel

	detect = packet.GetDetection()

	geo := packet.GetGeometry()
	if geo != nil {
		field = geo.GetField()
	}

	for _, robot := range detect.GetRobotsBlue() {
		x := float64(robot.GetX())
		y := float64(robot.GetY())
		w := float64(*robot.Orientation)

		wp.gamestate.SetRobot(robot.GetRobotId(), x, y, w, true)
	}

	for _, robot := range detect.GetRobotsYellow() {
		x := float64(robot.GetX())
		y := float64(robot.GetY())
		w := float64(*robot.Orientation)

		wp.gamestate.SetRobot(robot.GetRobotId(), x, y, w, false)

	}

	for _, ball := range detect.GetBalls() {
		x := float64(ball.GetX())
		y := float64(ball.GetY())
		z := float64(ball.GetZ())

		wp.gamestate.SetBall(x, y, z)
	}

	parseFieldData(&wp.gamestate.Field, field)
}

// Start a SSL Vision receiver, returns a channel from
// which SSL wrapper packets can be obtained.
func (wp *WorldPredictor) setupSSLVisionReceiver(addr string) {
	wp.ssl_receiver = receiver.NewSSLReceiver(addr)
	wp.ssl_receiver.Connect()

	wp.ssl_receiver_channel = make(chan ssl_vision.SSL_WrapperPacket)
	go wp.ssl_receiver.Receive(wp.ssl_receiver_channel)
}

func NewWorldPredictor(sslReceiverAddress string, gs *gamestate.GameState) *WorldPredictor {
	wp := &WorldPredictor{}
	wp.gamestate = gs
	wp.setupSSLVisionReceiver(sslReceiverAddress)
	return wp
}

// Parse geoemtry field data
func parseFieldData(f *gamestate.Field, data *ssl_vision.SSL_GeometryFieldSize) {
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
func parseFieldLines(f *gamestate.Field, lines []*ssl_vision.SSL_FieldLineSegment) {
	for _, line := range lines {
		if hasLine(line.GetName(), f) {
			continue
		}
		p1 := line.GetP1()
		p2 := line.GetP2()
		f.SetLine(
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
func parseFieldArcs(f *gamestate.Field, arcs []*ssl_vision.SSL_FieldCircularArc) {
	for _, arc := range arcs {
		if hasArc(arc.GetName(), f) {
			continue
		}

		center := arc.GetCenter()
		f.SetArc(
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
func convertShapeType(typ ssl_vision.SSL_FieldShapeType) gamestate.FieldShape {
	return gamestate.FieldShape(typ)
}

// Check if Field contains some line
// with given name.
func hasLine(name string, f *gamestate.Field) bool {
	for _, line := range f.FieldLines {
		if line.Name == name {
			return true
		}
	}

	return false
}

// Check if Field contains some arc
// with given name.
func hasArc(name string, f *gamestate.Field) bool {
	for _, arc := range f.FieldArcs {
		if arc.Name == name {
			return true
		}
	}

	return false
}
