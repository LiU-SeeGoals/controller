package examples

import (
	"math"
	"time"
	"net"
	"fmt"
	"strconv"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/client"
	"gonum.org/v1/gonum/mat"
	"google.golang.org/protobuf/proto"

	"github.com/LiU-SeeGoals/controller/internal/proto/basestation" // This import path assumes 'controller' is the module name

)


func main() {
	var port int = 25565
	go startServer(port)
	sendPacket(port)
}

func sendPacket(port int) {
	BaseStationClient := client.NewBaseStationClient("127.0.0.1:"+strconv.Itoa(port))
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
	time.Sleep(2 * time.Second)

	BaseStationClient.Send(actions) // Send the messages again for fun
	time.Sleep(2 * time.Second)

}

func startServer(port int) {
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
	for {
		_, remoteaddr, err := ser.ReadFromUDP(p)
		fmt.Printf("Read a message from %v \n", remoteaddr)

		command := &basestation.Command{}
		proto.Unmarshal(p, command)
		fmt.Printf("Unmarshalled command: %v\n", command)

		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
		
	}
}