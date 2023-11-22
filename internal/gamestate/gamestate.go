package gamestate

import (
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
	ball        *Ball
}

func (gs *GameState) Test() {
	act := &action.Move{}
	act.Pos = gs.yellow_team[0].pos
	act.Dest = mat.NewVecDense(3, nil)
	act.Dest.SetVec(0, 4)
	act.Dest.SetVec(1, 4)
	act.Dest.SetVec(2, 0)
	act.Dribble = true

	//act := &action.Kick{}
	//act.Kickspeed = 10

	//act := &action.Dribble{}
	//act.Dribble = true

	var action []action.Action
	action = append(action, act)

	gs.Grsim_client.AddActions(action)
}

// Updates position of robots and balls to their actual position
func (gs *GameState) Update() {
	var packet ssl_vision.SSL_WrapperPacket

	var detect *ssl_vision.SSL_DetectionFrame = nil
	packet = <-gs.ssl_receiver_channel

	detect = packet.GetDetection()

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
