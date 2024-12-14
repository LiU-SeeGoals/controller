package client

import (
	"fmt"
	"net"

	"github.com/LiU-SeeGoals/controller/internal/helper"
	"github.com/LiU-SeeGoals/controller/internal/state"
	"github.com/LiU-SeeGoals/proto_go/ssl_vision"
	"google.golang.org/protobuf/proto"
)

const (
	// Read buffer size
	READ_BUFFER_SIZE = 8192
)

// SSL Vision receiver
type SSLConnection struct {
	// Connection
	conn *net.UDPConn
	// UDP address
	addr *net.UDPAddr
	// Read buffer
	buff []byte
	// SSL lets not heap allocate this every time
	packet ssl_vision.SSL_WrapperPacket
}

// Create a new SSL vision receiver.
// Address should be <ip>:<port>
func NewSSLConnection(addr string) *SSLConnection {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}

	return &SSLConnection{
		conn: nil,
		addr: udpAddr,
		buff: make([]byte, READ_BUFFER_SIZE),
	}
}

// Connect/subscribe receiver to UDP multicast.
// Note, this will NOT block.
func (r *SSLConnection) Connect() {
	conn, err := net.ListenMulticastUDP("udp", nil, r.addr)
	if err != nil {
		panic(err)
	}

	r.conn = conn
}

// Start receiving packets.
// This function should be run in a goroutine:
//
//	go recv.Receive()
//
// Parsed packets are transferred using packetChan.
func (r *SSLConnection) Receive(packetChan chan *ssl_vision.SSL_WrapperPacket) {
	for {
		sz, err := r.conn.Read(r.buff)
		if err != nil {
			fmt.Printf("Unable to receive packet: %s", err)
			continue
		}

		err = proto.Unmarshal(r.buff[:sz], &r.packet)
		if err != nil {
			fmt.Printf("Unable to unmarshal packet: %s", err)
			continue
		}
		packetChan <- &r.packet
	}
}

type SSLVisionClient struct {
	ssl             *SSLConnection
	ssl_channel_in  chan *ssl_vision.SSL_WrapperPacket
	ssl_channel_out chan *ssl_vision.SSL_WrapperPacket
}

func unpack(packet *ssl_vision.SSL_WrapperPacket, gs *state.GameState, play_time int64) {
	detect := packet.GetDetection()
	gs.SetMessageReceivedTime(play_time)

	for _, robot := range detect.GetRobotsBlue() {
		x := robot.GetX()
		y := robot.GetY()
		angle := robot.GetOrientation()
		// fmt.Println("Robot", robot.GetRobotId(), "x:", x, "y:", y, "angle:", angle)

		gs.SetBlueRobot(robot.GetRobotId(), x, y, angle, play_time)
	}

	for _, robot := range detect.GetRobotsYellow() {
		// fmt.Println("Robot", robot.GetRobotId(), "x:", robot.GetX(), "y:", robot.GetY(), "angle:", robot.GetOrientation())
		x := robot.GetX()
		y := robot.GetY()
		angle := robot.GetOrientation()
		gs.SetYellowRobot(robot.GetRobotId(), x, y, angle, play_time)

	}

	for _, ball := range detect.GetBalls() {
		// fmt.Println("Ball", ball.GetX(), ball.GetY(), ball.GetZ())
		x := ball.GetX()
		y := ball.GetY()
		z := ball.GetZ()

		gs.SetBall(x, y, z, play_time)
	}
	gs.SetValid(true)

	geometry := packet.GetGeometry()
	field := geometry.GetField()

	gs.SetField(field.GetFieldLength(),
		field.GetFieldWidth(),
		field.GetGoalWidth(),
		field.GetGoalDepth(),
		field.GetBoundaryWidth(),
		field.GetPenaltyAreaDepth(),
		field.GetPenaltyAreaWidth(),
	)
	for _, line := range field.GetFieldLines() {
		gs.AddFieldLine(line.GetName(), line.GetP1().GetX(), line.GetP1().GetY(), line.GetP2().GetX(), line.GetP2().GetY(), line.GetThickness(), int(line.GetType()))
	}
	for _, arc := range field.GetFieldArcs() {
		gs.AddFieldArc(arc.GetName(), arc.GetCenter().GetX(), arc.GetCenter().GetY(), arc.GetRadius(), arc.GetA1(), arc.GetA2(), arc.GetThickness(), int(arc.GetType()))
	}

}

func (receiver *SSLVisionClient) InitGameState(gs *state.GameState, play_time int64) {
	packet, ok := <-receiver.ssl_channel_out
	if !ok {
		fmt.Println("SSL Channel closed")
		return
	}
	unpack(packet, gs, play_time)
}

func (receiver *SSLVisionClient) UpdateGamestate(gs *state.GameState, play_time int64) {
	packet, ok := <-receiver.ssl_channel_out

	if !ok {
		fmt.Println("SSL Channel closed")
		return
	}
	unpack(packet, gs, play_time)
}

// Start a SSL Vision receiver, returns a channel from
// which SSL wrapper packets can be obtained.
func (receiver *SSLVisionClient) Connect() {
	receiver.ssl.Connect()
	go receiver.ssl.Receive(receiver.ssl_channel_in)
}

func NewSSLVisionClient(sslReceiverAddress string) *SSLVisionClient {
	ssl_channel_in, ssl_channel_out := helper.NB_KeepLatestChan[*ssl_vision.SSL_WrapperPacket]()
	receiver := &SSLVisionClient{
		ssl:             NewSSLConnection(sslReceiverAddress),
		ssl_channel_in:  ssl_channel_in,
		ssl_channel_out: ssl_channel_out,
	}
	receiver.Connect()
	return receiver
}
