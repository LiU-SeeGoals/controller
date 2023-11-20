package gamestate

import (
	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"github.com/LiU-SeeGoals/controller/internal/proto/ssl_vision"
	"github.com/LiU-SeeGoals/controller/internal/receiver"
)

const TEAM_SIZE = 6

type GameState struct {
	Grsim_client         *client.GrsimClient
	ssl_receiver         *receiver.SSLReceiver
	ssl_receiver_channel chan ssl_vision.SSL_WrapperPacket

	blue_team   [TEAM_SIZE]*Robot
	yellow_team [TEAM_SIZE]*Robot
	ball        *Ball
	field       Field
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

	gs.Grsim_client = client.NewSSLGrsimClient(sslClientAddress)
	gs.Grsim_client.Connect()

	gs.setupSSLVisionReceiver(sslReceiverAddress)

	gs.ball = NewBall()

	for i := 0; i < TEAM_SIZE; i++ {
		gs.blue_team[i] = NewRobot(i, Blue)
		gs.yellow_team[i] = NewRobot(i, Yellow)
	}

	return gs
}

func (gs *GameState) String() string {
	gs_str := "{\n blue team: {\n"
	for i := 0; i < TEAM_SIZE; i++ {
		gs_str += "{" + gs.blue_team[i].String() + " },\n"
	}
	gs_str += "},\n yellow team: {\n"
	for i := 0; i < TEAM_SIZE; i++ {
		gs_str += "{" + gs.yellow_team[i].String() + " },\n"
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

func parseFieldArcs(f *Field, lines []*ssl_vision.SSL_FieldCircularArc) {

}

// Glorified type cast
// Converts ssl vision enum to our own enum
func convertShapeType(typ ssl_vision.SSL_FieldShapeType) gamestate.FieldShape {
	return typ
}
