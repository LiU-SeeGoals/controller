package client

import (
	"fmt"
	"math"
	"net"
	"strconv"
	"testing"

	"gonum.org/v1/gonum/mat"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/proto-messages/robot_action"
	"google.golang.org/protobuf/proto"
)

var globalCommand *robot_action.Command

// This test starts a client and a sever and then sends a action to the server
// and then checks if the response matches what was sent.
// Only checks commandID and robotID in the message.
func TestSocketCommunication(t *testing.T) {

	// Define stop action
	stopAction := &action.Stop{Id: 2}
	stopCommand := &robot_action.Command{CommandId: robot_action.ActionType_STOP_ACTION, RobotId: 2}

	// Define kick action
	kickAction := &action.Kick{Id: 6, KickSpeed: 5}
	kickCommand := &robot_action.Command{CommandId: robot_action.ActionType_KICK_ACTION, RobotId: 6, KickSpeed: 5}

	// Define init action
	initAction := &action.Init{Id: 3}
	initCommand := &robot_action.Command{CommandId: robot_action.ActionType_INIT_ACTION, RobotId: 3}

	// Define move action
	moveAction := &action.Move{
		Id:   1,
		Pos:  mat.NewVecDense(3, []float64{100, 200, math.Pi}),
		Dest: mat.NewVecDense(3, []float64{300, 400, -math.Pi}),
	}
	moveCommand := &robot_action.Command{
		CommandId: robot_action.ActionType_MOVE_ACTION,
		RobotId:   1,
		Pos:       &robot_action.Vector3D{X: int32(100), Y: int32(200), W: float32(math.Pi)},
		Dest:      &robot_action.Vector3D{X: int32(300), Y: int32(400), W: float32(-math.Pi)},
	}

	// Define set navigation direction action.
	setNavDirAction := &action.SetNavigationDirection{
		Id:        9,
		Direction: mat.NewVecDense(2, []float64{100, 200}),
	}
	setNavDirCommand := &robot_action.Command{
		CommandId: robot_action.ActionType_SET_NAVIGATION_DIRECTION_ACTION,
		RobotId:   9,
		Direction: &robot_action.Vector2D{X: int32(100), Y: int32(200)},
	}

	// Define rotate.
	rotateAction := &action.Rotate{
		Id:         3,
		AngularVel: 5,
	}
	rotateCommand := &robot_action.Command{
		CommandId:  robot_action.ActionType_ROTATE_ACTION,
		RobotId:    3,
		AngularVel: 5,
	}

	// Test cases
	testCases := []struct {
		input    action.Action
		expected *robot_action.Command
	}{
		{stopAction, stopCommand},
		{kickAction, kickCommand},
		{initAction, initCommand},
		{moveAction, moveCommand},
		{setNavDirAction, setNavDirCommand},
		{rotateAction, rotateCommand},
	}

	commandChan := make(chan *robot_action.Command)
	var port int = 25565
	go startServer(port, commandChan)

	for _, tc := range testCases {

		command := testCommunication(tc.input, commandChan, port)

		if command.GetRobotId() != tc.expected.GetRobotId() {
			t.Errorf("Expected: %v, got: %v", tc.expected, command)
		}

		if command.GetCommandId() != tc.expected.GetCommandId() {
			t.Errorf("Expected: %v, got: %v", tc.expected, command)
		}
	}
}

func testCommunication(newCommand action.Action, commandChan chan *robot_action.Command, port int) *robot_action.Command {

	BaseStationClient := NewBaseStationClient("127.0.0.1:" + strconv.Itoa(port))
	BaseStationClient.Init()
	BaseStationClient.SendActions([]action.Action{newCommand})

	command := <-commandChan

	return command
}

func startServer(port int, commandChan chan<- *robot_action.Command) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("127.0.0.1"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
	for {
		p := make([]byte, 32) // Reinitialize before each read

		ser.ReadFromUDP(p)
		command := &robot_action.Command{}
		proto.Unmarshal(p, command)

		// Send the command to the channel
		commandChan <- command
	}
}
