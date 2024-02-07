package world_predictor

import (
	"github.com/LiU-SeeGoals/controller/internal/config"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"github.com/LiU-SeeGoals/proto-messages/go/ssl_vision"
	"gonum.org/v1/gonum/mat"

	"github.com/LiU-SeeGoals/controller/internal/receiver"
)

type WorldPredictor struct {
	ssl_receiver         *receiver.SSLReceiver
	buffer *DoubleRingBuffer
}

func NewWorldPredictor() *WorldPredictor {
	wp := &WorldPredictor{}
	wp.ssl_receiver = receiver.NewSSLReceiver(config.GetSSLClientAddress())
	wp.ssl_receiver.Connect()
	wp.buffer = NewDoubleRingBuffer()
	go wp.predictGameState()
	return wp
}

func (wp *WorldPredictor) GetGameState(i int) *gamestate.GameState {
	return wp.buffer.GetGameState(i)
}

func (wp *WorldPredictor) Update() {
	wp.buffer.Update()
}

func (wp *WorldPredictor) predictGameState() {
	var packet *ssl_vision.SSL_WrapperPacket
	var curGameState *gamestate.GameState = wp.buffer.GetGameStateInProgress()
	var curFrameNumber uint32 = 0
	var amountOfPackets uint32 = 0
	robotNormalizationFactor := make(map[gamestate.Team]map[int]float64)
	robotNormalizationFactor[gamestate.Blue] = make(map[int]float64)
	robotNormalizationFactor[gamestate.Yellow] = make(map[int]float64)
	var ballNormalizationFactor float64 = 0

	for {
		packet = wp.ssl_receiver.Receive()

		if wp.isGameStateDone(*packet.Detection.FrameNumber, curFrameNumber, amountOfPackets) {
			// Normalize game state and place it in buffer.
			wp.normalizeGameState(curGameState, robotNormalizationFactor, &ballNormalizationFactor)
			wp.buffer.PlaceGameState()

			// Prepare to handle next gamestate.
			curGameState = wp.buffer.GetGameStateInProgress()
			curFrameNumber = *packet.Detection.FrameNumber
			robotNormalizationFactor = make(map[gamestate.Team]map[int]float64)
			robotNormalizationFactor[gamestate.Blue] = make(map[int]float64)
			robotNormalizationFactor[gamestate.Yellow] = make(map[int]float64)
			ballNormalizationFactor = 0
			amountOfPackets = 0
		}

		wp.predictRobot(packet, curGameState, robotNormalizationFactor)
		wp.predictBall(packet, curGameState, &ballNormalizationFactor)
		amountOfPackets++
	}
}

func (wp *WorldPredictor) isGameStateDone(packetNumber, frameNumber, amountOfPacketsReceived uint32) bool {
	isCurrentFramePacket := packetNumber == frameNumber
	isAllFramePacketReceived := amountOfPacketsReceived >= config.GetAmountOfCameras()
	return !isCurrentFramePacket || isAllFramePacketReceived
}

func (wp *WorldPredictor) predictRobot(packet *ssl_vision.SSL_WrapperPacket, g *gamestate.GameState, robotNormalizationFactor map[gamestate.Team]map[int]float64) {
	blue_team := packet.Detection.RobotsBlue
	yellow_team := packet.Detection.RobotsYellow
	
	var robot *gamestate.Robot
	var robotId int
	var robotPosition *mat.VecDense
	var confidence float64
	var curNormFactor float64
	var keyExist bool
	var x, y, w float64

	for _, robotDetection := range blue_team {
		robotId = int(robotDetection.GetRobotId())

		robot = g.GetRobot(robotId, true)
		if robot == nil {
			g.AddRobot(robotId, true)
			robot = g.GetRobot(robotId, true)
		}

		confidence = float64(robotDetection.GetConfidence())
		curNormFactor, keyExist = robotNormalizationFactor[gamestate.Blue][robotId]
		if !keyExist {
			robotNormalizationFactor[gamestate.Blue][robotId] = 0
		}
		robotNormalizationFactor[gamestate.Blue][robotId] = curNormFactor + confidence

		robotPosition = robot.GetPosition()
		x = robotPosition.AtVec(0) + float64(robotDetection.GetX()) * confidence
		y = robotPosition.AtVec(1) + float64(robotDetection.GetY()) * confidence
		w = robotPosition.AtVec(2) + float64(robotDetection.GetOrientation()) * confidence
		robot.SetPosition(x, y, w)

		// Add speed calculation here
		
	}

	for _, robotDetection := range yellow_team {
		robotId = int(robotDetection.GetRobotId())

		robot = g.GetRobot(robotId, false)
		if robot == nil {
			g.AddRobot(robotId, false)
			robot = g.GetRobot(robotId, false)
		}

		confidence = float64(robotDetection.GetConfidence())
		curNormFactor, keyExist = robotNormalizationFactor[gamestate.Yellow][robotId] 
		if !keyExist {
			robotNormalizationFactor[gamestate.Yellow][robotId] = 0
		}
		robotNormalizationFactor[gamestate.Yellow][robotId] = curNormFactor + confidence

		robotPosition = robot.GetPosition()
		x = robotPosition.AtVec(0) + float64(robotDetection.GetX()) * confidence
		y = robotPosition.AtVec(1) + float64(robotDetection.GetY()) * confidence
		w = robotPosition.AtVec(2) + float64(robotDetection.GetOrientation()) * confidence
		robot.SetPosition(x, y, w)

		// Add speed calculation here
	}
}

func (wp *WorldPredictor) predictBall(packet *ssl_vision.SSL_WrapperPacket, gamestate *gamestate.GameState, ballNormalizationFactor *float64) {
	balls := packet.Detection.Balls
	for _, ball := range balls {
		confidence := float64(ball.GetConfidence())
		*ballNormalizationFactor += confidence
		ballPosition := gamestate.GetBall().GetPosition()
		x := ballPosition.AtVec(0) + float64(ball.GetX())
		y := ballPosition.AtVec(1) + float64(ball.GetY())
		z := ballPosition.AtVec(2) + float64(ball.GetZ())
		gamestate.SetBall(x, y, z)
	}
}

func (wp *WorldPredictor) normalizeGameState(g *gamestate.GameState, robotNormalizationFactor map[gamestate.Team]map[int]float64, ballNormalizationFactor *float64) {
	var robot *gamestate.Robot
	for robotId, _ := range robotNormalizationFactor[gamestate.Blue] {
		robot = g.GetRobot(robotId, true)
		robot.NormalizePosition(robotNormalizationFactor[gamestate.Blue][robotId])
	}
	for robotId, _ := range robotNormalizationFactor[gamestate.Yellow] {
		robot = g.GetRobot(robotId, false)
		robot.NormalizePosition(robotNormalizationFactor[gamestate.Yellow][robotId])
	}

	g.GetBall().NormalizePosition(*ballNormalizationFactor)
}

// Updates position of robots and balls to their actual position
// func (wp *WorldPredictor) Update() {
// 	var packet ssl_vision.SSL_WrapperPacket
// 
// 	var detect *ssl_vision.SSL_DetectionFrame
// 	var field *ssl_vision.SSL_GeometryFieldSize
// 
// 	packet = <-wp.ssl_receiver_channel
// 
// 	detect = packet.GetDetection()
// 
// 	geo := packet.GetGeometry()
// 	if geo != nil {
// 		field = geo.GetField()
// 	}
// 
// 	for _, robot := range detect.GetRobotsBlue() {
// 		x := float64(robot.GetX())
// 		y := float64(robot.GetY())
// 		w := float64(*robot.Orientation)
// 
// 		wp.gamestate.SetRobot(robot.GetRobotId(), x, y, w, true)
// 	}
// 
// 	for _, robot := range detect.GetRobotsYellow() {
// 		x := float64(robot.GetX())
// 		y := float64(robot.GetY())
// 		w := float64(*robot.Orientation)
// 
// 		wp.gamestate.SetRobot(robot.GetRobotId(), x, y, w, false)
// 
// 	}
// 
// 	for _, ball := range detect.GetBalls() {
// 		x := float64(ball.GetX())
// 		y := float64(ball.GetY())
// 		z := float64(ball.GetZ())
// 
// 		wp.gamestate.SetBall(x, y, z)
// 	}
// 
// 	parseFieldData(&wp.gamestate.Field, field)
// }

// Start a SSL Vision receiver, returns a channel from
// which SSL wrapper packets can be obtained.
// func (wp *WorldPredictor) setupSSLVisionReceiver(addr string) {
// 	wp.ssl_receiver = receiver.NewSSLReceiver(addr)
// 	wp.ssl_receiver.Connect()
// 
// 	wp.ssl_receiver_channel = make(chan ssl_vision.SSL_WrapperPacket)
// 	go wp.ssl_receiver.Receive(wp.ssl_receiver_channel)
// }
// 
// func NewWorldPredictor() *WorldPredictor {
// 	wp := &WorldPredictor{}
// 	wp.gamestate = gs
// 	wp.setupSSLVisionReceiver(sslReceiverAddress)
// 	return wp
// }
// 
// // Parse geoemtry field data
// func parseFieldData(f *gamestate.Field, data *ssl_vision.SSL_GeometryFieldSize) {
// 	if data == nil {
// 		return
// 	}
// 
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
// 
// 	parseFieldLines(f, data.GetFieldLines())
// 	parseFieldArcs(f, data.GetFieldArcs())
// }
// 
// // Parse field lines from ssl packet
// //
// // Field object should be passed from game state object.
// func parseFieldLines(f *gamestate.Field, lines []*ssl_vision.SSL_FieldLineSegment) {
// 	for _, line := range lines {
// 		if hasLine(line.GetName(), f) {
// 			continue
// 		}
// 		p1 := line.GetP1()
// 		p2 := line.GetP2()
// 		f.SetLine(
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
// 
// // Parse arcs from ssl packet
// //
// // Field object should be passed from game state object.
// func parseFieldArcs(f *gamestate.Field, arcs []*ssl_vision.SSL_FieldCircularArc) {
// 	for _, arc := range arcs {
// 		if hasArc(arc.GetName(), f) {
// 			continue
// 		}
// 
// 		center := arc.GetCenter()
// 		f.SetArc(
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
// 
// // Glorified type cast
// // Converts ssl vision enum to our own enum
// func convertShapeType(typ ssl_vision.SSL_FieldShapeType) gamestate.FieldShape {
// 	return gamestate.FieldShape(typ)
// }
// 
// // Check if Field contains some line
// // with given name.
// func hasLine(name string, f *gamestate.Field) bool {
// 	for _, line := range f.FieldLines {
// 		if line.Name == name {
// 			return true
// 		}
// 	}
// 
// 	return false
// }
// 
// // Check if Field contains some arc
// // with given name.
// func hasArc(name string, f *gamestate.Field) bool {
// 	for _, arc := range f.FieldArcs {
// 		if arc.Name == name {
// 			return true
// 		}
// 	}
// 
// 	return false
// }
