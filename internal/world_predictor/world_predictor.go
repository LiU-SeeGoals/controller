package world_predictor

import (
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"github.com/LiU-SeeGoals/controller/internal/receiver"
	"github.com/LiU-SeeGoals/proto_go/ssl_vision"
)

type WorldPredictor struct {
	ssl_receiver         *receiver.SSLReceiver
	ssl_receiver_channel chan ssl_vision.SSL_WrapperPacket
	gamestate            *gamestate.GameState
}

// Update the gamestate with new information derived from the latest SSL Vision packet. 
func (wp *WorldPredictor) UpdateGamestate() {
	var packet ssl_vision.SSL_WrapperPacket
	var detect *ssl_vision.SSL_DetectionFrame

	packet = <-wp.ssl_receiver_channel

	detect = packet.GetDetection()

	for _, robot := range detect.GetRobotsBlue() {
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

}

// Start a SSL Vision receiver, returns a channel from
// which SSL wrapper packets can be obtained.
func (wp *WorldPredictor) setupSSLVisionReceiver(addr string) {
	wp.ssl_receiver = receiver.NewSSLReceiver(addr)
	wp.ssl_receiver.Connect()

	wp.ssl_receiver_channel = make(chan ssl_vision.SSL_WrapperPacket)
	go wp.ssl_receiver.Receive(wp.ssl_receiver_channel)
}

// Constructor for WorldPredictor
func NewWorldPredictor(sslReceiverAddress string, gs *gamestate.GameState) *WorldPredictor {
	wp := &WorldPredictor{}
	wp.gamestate = gs
	wp.setupSSLVisionReceiver(sslReceiverAddress)
	return wp
}

