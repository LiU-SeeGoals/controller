package receiver

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

type SSLReceiver struct {
	ssl         *SSLConnection
	ssl_channel chan *ssl_vision.SSL_WrapperPacket
}

type GameState interface {
	SetYellowRobot(robotId uint32, x, y, w float64)
	SetBlueRobot(robotId uint32, x, y, w float64)
	SetBall(x, y, z float64)
}

func (receiver *SSLReceiver) UpdateGamestate(gs GameState) {
	packet, ok := <-receiver.ssl_channel
	if !ok {
		return
	}
	detect := packet.GetDetection()

	for _, robot := range detect.GetRobotsBlue() {
		x := float64(robot.GetX())
		y := float64(robot.GetY())
		w := float64(*robot.Orientation)

		gs.SetYellowRobot(robot.GetRobotId(), x, y, w)
	}

	for _, robot := range detect.GetRobotsYellow() {
		x := float64(robot.GetX())
		y := float64(robot.GetY())
		w := float64(*robot.Orientation)

		gs.SetBlueRobot(robot.GetRobotId(), x, y, w)
	}

	for _, ball := range detect.GetBalls() {
		x := float64(ball.GetX())
		y := float64(ball.GetY())
		z := float64(ball.GetZ())

		gs.SetBall(x, y, z)
	}

}

// Start a SSL Vision receiver, returns a channel from
// which SSL wrapper packets can be obtained.
func (receiver *SSLReceiver) Connect() {
	receiver.ssl.Connect()
	go receiver.ssl.Receive(receiver.ssl_channel)
}

func NewSSLReceiver(sslReceiverAddress string) *SSLReceiver {
	receiver := &SSLReceiver{
		ssl:         NewSSLConnection(sslReceiverAddress),
		ssl_channel: make(chan *ssl_vision.SSL_WrapperPacket),
	}
	receiver.Connect()
	return receiver
}
