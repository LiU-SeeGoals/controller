package client

import (
	"time"
	"net"
	"fmt"
	"strconv"
	"testing"
    "math"
    "sync"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"google.golang.org/protobuf/proto"
    "gonum.org/v1/gonum/mat"

	"github.com/LiU-SeeGoals/controller/internal/proto/basestation" // This import path assumes 'controller' is the module name

)

type Response struct {
	Message *basestation.Command
}

func TestSendAction(t *testing.T) {
    var wg sync.WaitGroup
    wg.Add(1)
    var port int = 25565
    responseChan := make(chan Response)
	go startServer(port, &wg, responseChan)
	sendPacket(port)

    
    
    wg.Wait()
    close(responseChan)

    fmt.Println("checking response")
    for response := range responseChan {
		command := response.Message
		fmt.Printf("Received command: %+v\n", command)
	}
    fmt.Println("checking response2")

}


func sendPacket(port int) {
	BaseStationClient := NewBaseStationClient("127.0.0.1:"+strconv.Itoa(port))
	BaseStationClient.Init()

	// Creates 2 random actions to send
	actions := []action.Action{
		&action.Stop{Id: 2},
		&action.Move{
			Id: 3,
			Pos: mat.NewVecDense(3, []float64{100, 200, math.Pi}), // Example values for Pos
			Goal: mat.NewVecDense(3, []float64{300, 400, -math.Pi}), // Example values for Goal
		},
	}

	BaseStationClient.Send(actions)
	// time.Sleep(2 * time.Second)

	// BaseStationClient.Send(actions) // Send the messages again for fun
	// time.Sleep(2 * time.Second)

}

func startServer(port int, wg *sync.WaitGroup, actionChan chan<- Response) {
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
    ser.SetReadDeadline(time.Now().Add(1 * time.Second))
    ser.ReadFromUDP(p)
    //fmt.Printf("Read a message from %v \n", remoteaddr)

    command := &basestation.Command{}
    proto.Unmarshal(p, command)
    fmt.Printf("Unmarshalled command: %v\n", command)

    response := Response{Message: command}
    actionChan <- response
    
		
	
}
