package client

import (
	"fmt"
	"net"

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
	ssl         *SSLConnection
	ssl_channel chan *ssl_vision.SSL_WrapperPacket
}

type GameState interface {
	SetYellowRobot(robotId uint32, x, y, w float64, time int64)
	SetBlueRobot(robotId uint32, x, y, w float64, time int64)
	SetBall(x, y, z float64, time int64)
	SetMessageReceivedTime(time int64)
	SetLagTime(lagTime int64)
	GetMessageReceivedTime() int64
}

func unpack(packet *ssl_vision.SSL_WrapperPacket, gs GameState, play_time int64) {
	detect := packet.GetDetection()
	gs.SetMessageReceivedTime(play_time)
	gs.SetLagTime(0)

	for _, robot := range detect.GetRobotsBlue() {
		x := float64(robot.GetX())
		y := float64(robot.GetY())
		w := float64(*robot.Orientation)

		gs.SetYellowRobot(robot.GetRobotId(), x, y, w, play_time)
	}

	for _, robot := range detect.GetRobotsYellow() {
		x := float64(robot.GetX())
		y := float64(robot.GetY())
		w := float64(*robot.Orientation)

		gs.SetBlueRobot(robot.GetRobotId(), x, y, w, play_time)
	}

	for _, ball := range detect.GetBalls() {
		x := float64(ball.GetX())
		y := float64(ball.GetY())
		z := float64(ball.GetZ())

		gs.SetBall(x, y, z, play_time)
	}
}

func update_lag_time(gs GameState, play_time int64) {
	time := gs.GetMessageReceivedTime()
	lag_time := play_time - time
	gs.SetLagTime(lag_time)
}

func (receiver *SSLVisionClient) InitGameState(gs GameState, play_time int64) {
	packet, ok := <-receiver.ssl_channel
	if !ok {
		fmt.Println("SSL Channel closed")
		return
	}
	unpack(packet, gs, play_time)
}

func (receiver *SSLVisionClient) UpdateGamestate(gs GameState, play_time int64) {
	packet, ok := <-receiver.ssl_channel
	if !ok {
		fmt.Println("SSL Channel closed")
		return
	}
	unpack(packet, gs, play_time)
}

func (receiver *SSLVisionClient) UpdateGamestateNB(gs GameState, play_time int64) {
	// // none-blocking receive
	select {
	case packet := <-receiver.ssl_channel:
		unpack(packet, gs, play_time)
	default:
		update_lag_time(gs, play_time)
	}

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
