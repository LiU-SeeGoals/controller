package client

import (
	"time"
	"net"
	"fmt"
	"strconv"
	"testing"
    "sync"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"google.golang.org/protobuf/proto"
	"github.com/LiU-SeeGoals/controller/internal/proto/basestation"
)

var globalCommand *basestation.Command

// This test starts a client and a sever and then sends a action to the server
// and then checks if the response matches what was sent.
// Only checks commandID and robotID in the message.
func TestSocketCommunication(t *testing.T) {
	testCases := []struct {
		input         action.Action
		expected      *basestation.Command
	}{
		{&action.Stop{Id: 2}, &basestation.Command{CommandId: basestation.ActionType_STOP_ACTION,RobotId: 2,}},
		{&action.Kick{Id: 6, Speed: 5}, &basestation.Command{CommandId: basestation.ActionType_KICK_ACTION,RobotId: 6, Speed: 5,}},
	}


	for _, tc := range testCases {
		command := testCommunication(tc.input)
		time.Sleep(1000 * time.Millisecond)

		if command.GetRobotId() != tc.expected.GetRobotId() {
			t.Errorf("Expected: %v, got: %v", tc.expected, command)
		}

		if command.GetCommandId() != tc.expected.GetCommandId() {
			t.Errorf("Expected: %v, got: %v", tc.expected, command)
		}
	}
}

func testCommunication(newCommand action.Action) *basestation.Command{
	var wg sync.WaitGroup
    wg.Add(1)
    var port int = 25565

	go startServer(port, &wg)

	BaseStationClient := NewBaseStationClient("127.0.0.1:"+strconv.Itoa(port))
	BaseStationClient.Init()
	BaseStationClient.Send([]action.Action{newCommand})

    wg.Wait()
	return globalCommand
	
	
}

func startServer(port int, wg *sync.WaitGroup) {
    defer wg.Done()
	p := make([]byte, 32)
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("127.0.0.1"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
		return
	}
    ser.SetReadDeadline(time.Now().Add(1 * time.Second)) // one sec timeout
    ser.ReadFromUDP(p)
    command := &basestation.Command{}
    proto.Unmarshal(p, command)

    globalCommand = command
	ser.Close()
		
	
}
