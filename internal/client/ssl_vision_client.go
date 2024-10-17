package client

import (
	"fmt"
	"net"

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
		// Clear the channel if something is there
		fmt.Println("Received packet")
		fmt.Println("size", len(packetChan))
		if len(packetChan) > 0 {
			select {
			case <-packetChan:
			default:
			}
		}
		packetChan <- &r.packet
	}
}

type SSLVisionClient struct {
	ssl         *SSLConnection
	ssl_channel chan *ssl_vision.SSL_WrapperPacket
}

func unpack(packet *ssl_vision.SSL_WrapperPacket, gs *state.GameState, play_time int64) {
	detect := packet.GetDetection()
	gs.SetMessageReceivedTime(play_time)

	for _, robot := range detect.GetRobotsBlue() {
		x := robot.GetX()
		y := robot.GetY()
		angle := robot.GetOrientation()
		fmt.Println("Robot", robot.GetRobotId(), "x:", x, "y:", y, "angle:", angle)

		gs.SetBlueRobot(robot.GetRobotId(), x, y, angle, play_time)
	}

	for _, robot := range detect.GetRobotsYellow() {
		x := robot.GetX()
		y := robot.GetY()
		angle := robot.GetOrientation()
		gs.SetYellowRobot(robot.GetRobotId(), x, y, angle, play_time)

	}

	for _, ball := range detect.GetBalls() {
		x := ball.GetX()
		y := ball.GetY()
		z := ball.GetZ()

		gs.SetBall(x, y, z, play_time)
	}
	gs.SetValid(true)
}

func (receiver *SSLVisionClient) InitGameState(gs *state.GameState, play_time int64) {
	packet, ok := <-receiver.ssl_channel
	if !ok {
		fmt.Println("SSL Channel closed")
		return
	}
	unpack(packet, gs, play_time)
}

func (receiver *SSLVisionClient) UpdateGamestate(gs *state.GameState, play_time int64) {
	packet, ok := <-receiver.ssl_channel
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
	go receiver.ssl.Receive(receiver.ssl_channel)
}

func NewSSLVisionClient(sslReceiverAddress string) *SSLVisionClient {
	receiver := &SSLVisionClient{
		ssl:         NewSSLConnection(sslReceiverAddress),
		ssl_channel: make(chan *ssl_vision.SSL_WrapperPacket),
	}
	receiver.Connect()
	return receiver
}
