package examples

import (
	"fmt"
	"math"
	"net"
	"strconv"
	"time"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/proto-messages/robot_action"
	"gonum.org/v1/gonum/mat"
	"google.golang.org/protobuf/proto"
)

func main() {
	var port int = 25565
	go startServer(port)
	sendPacket(port)
}

func sendPacket(port int) {
	BaseStationClient := client.NewBaseStationClient("127.0.0.1:" + strconv.Itoa(port))
	BaseStationClient.Init()

	// Creates 2 random actions to send
	actions := []action.Action{
		&action.Stop{Id: 2},
		&action.Move{
			Id:   3,
			Pos:  mat.NewVecDense(3, []float64{100, 200, math.Pi}),
			Dest: mat.NewVecDense(3, []float64{300, 400, -math.Pi}),
		},
	}

	BaseStationClient.SendActions(actions)
	time.Sleep(2 * time.Second)

	BaseStationClient.SendActions(actions) // Send the messages again for fun
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

		command := &robot_action.Command{}
		proto.Unmarshal(p, command)
		fmt.Printf("Unmarshalled command: %v\n", command)

		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}

	}
}
