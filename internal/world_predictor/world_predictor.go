package world_predictor

import (
	"fmt"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"github.com/LiU-SeeGoals/controller/internal/receiver"
	"github.com/LiU-SeeGoals/controller/internal/webserver"
	"github.com/LiU-SeeGoals/proto_go/ssl_vision"
)

type WorldPredictor struct {
	ssl_receiver         *receiver.SSLReceiver
	ssl_receiver_channel chan ssl_vision.SSL_WrapperPacket
	gamestate            *gamestate.GameState
}

// Updates position of robots and balls to their actual position
func (wp *WorldPredictor) Update() {
	fmt.Println("Updating world predictor")
	var packet ssl_vision.SSL_WrapperPacket
	fmt.Println(packet)
	var detect *ssl_vision.SSL_DetectionFrame
	var field *ssl_vision.SSL_GeometryFieldSize

	packet = <-wp.ssl_receiver_channel

	detect = packet.GetDetection()

	geo := packet.GetGeometry()
	if geo != nil {
		fmt.Println("Got geometry")
		field = geo.GetField()
	}

	for _, robot := range detect.GetRobotsBlue() {
		fmt.Println("Robot: ", robot.GetRobotId())
		x := float64(robot.GetX())
		y := float64(robot.GetY())
		w := float64(*robot.Orientation)

		wp.gamestate.GetRobot(int(robot.GetRobotId()), gamestate.Blue).SetPosition(x, y, w)
	}

	for _, robot := range detect.GetRobotsYellow() {
		x := float64(robot.GetX())
		y := float64(robot.GetY())
		w := float64(*robot.Orientation)

		wp.gamestate.GetRobot(int(robot.GetRobotId()), gamestate.Yellow).SetPosition(x, y, w)

	}

	for _, ball := range detect.GetBalls() {
		x := float64(ball.GetX())
		y := float64(ball.GetY())
		z := float64(ball.GetZ())

		wp.gamestate.GetBall().SetPosition(x, y, z)
	}

	parseFieldData(wp.gamestate.Field, field)
	// wp.broadcastGameState()
	//wp.sendActions()
}

func (gs *WorldPredictor) handleIncoming(incomming []action.ActionDTO) {
	fmt.Println("Received a new action (gamestate)")

	// TODO also set manual control for the robot that is controlled

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

func (wp *WorldPredictor) broadcastGameState() {
	webserver.BroadcastGameState(wp.gamestate.ToJson())
	// list of incoming actions
	incomming := webserver.GetIncoming()

	// If we got new actions --> then handle them
	if len(incomming) > 0 {
		wp.handleIncoming(incomming)
	}
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

	f.FieldLength = data.GetFieldLength()
	f.FieldWidth = data.GetFieldWidth()
	fmt.Println("Field length: ", f.FieldLength)
	f.BoundaryWidth = data.GetBoundaryWidth()
	f.GoalDepth = data.GetGoalDepth()
	f.GoalWidth = data.GetGoalWidth()
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
