package client

import (
	"net"
	"time"

	"github.com/LiU-SeeGoals/controller/internal/helper"
	"github.com/LiU-SeeGoals/controller/internal/info"
	. "github.com/LiU-SeeGoals/controller/internal/logger"
	"github.com/LiU-SeeGoals/proto_go/ssl_vision"
	"google.golang.org/protobuf/proto"
)

const (
	// Read buffer size
	READ_BUFFER_SIZE = 8192
)

// SSL Vision receiver
type visionConnection struct {
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
func NewVisionConnection(addr string) *visionConnection {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}

	return &visionConnection{
		conn: nil,
		addr: udpAddr,
		buff: make([]byte, READ_BUFFER_SIZE),
	}
}

// Connect/subscribe receiver to UDP multicast.
// Note, this will NOT block.
func (r *visionConnection) Connect() {
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
func (r *visionConnection) Receive(packetChan chan *ssl_vision.SSL_WrapperPacket) {
	for {
		sz, err := r.conn.Read(r.buff)
		if err != nil {
			Logger.Errorf("Unable to receive packet: %v", err)
			continue
		}

		err = proto.Unmarshal(r.buff[:sz], &r.packet)
		if err != nil {
			Logger.Errorf("Unable to unmarshal packet: %v", err)
			continue
		}

		helper.NB_Send[ssl_vision.SSL_WrapperPacket](packetChan, &r.packet)
	}
}

type SSLVisionSource struct {
	ssl         *visionConnection
	ssl_channel chan *ssl_vision.SSL_WrapperPacket
}

func (receiver *SSLVisionSource) GetVisionData() (*ssl_vision.SSL_WrapperPacket) {
	select {
	case packet, ok := <-receiver.ssl_channel:
		if !ok {
			Logger.Warn("SSL Channel closed")
			return packet
		}

		return packet
	default:
		Logger.Warn("SSL Channel empty")
		return nil
	}
}

// Start a SSL Vision receiver, returns a channel from
// which SSL wrapper packets can be obtained.
func (receiver *SSLVisionSource) Connect() {
	receiver.ssl.Connect()
	go receiver.ssl.Receive(receiver.ssl_channel)
}

func NewSSLVisionSource(sslReceiverAddress string) *SSLVisionSource {
	ssl_channel := make(chan *ssl_vision.SSL_WrapperPacket)
	receiver := &SSLVisionSource{
		ssl:         NewVisionConnection(sslReceiverAddress),
		ssl_channel: ssl_channel,
	}
	receiver.Connect()
	return receiver
}

type DetectionBundle struct {
	// Frame metadata
	FrameNumber       uint32
	CameraID          uint32
	CaptureTime       float64
	SentTime          float64
	CaptureCameraTime float64

	// All raw detections
	RawBalls  []*ssl_vision.SSL_DetectionBall
	RawYellow []*ssl_vision.SSL_DetectionRobot
	RawBlue   []*ssl_vision.SSL_DetectionRobot

	// Lookup map: team → id → detection
	Robots map[info.Team]map[info.ID]*ssl_vision.SSL_DetectionRobot

	// ReceivedTime (local reception time, in ms)
	ReceivedTime int64
}

// PackDetectionFrame creates a DetectionBundle from a SSL_DetectionFrame.
func PackDetectionFrame(frame *ssl_vision.SSL_DetectionFrame) *DetectionBundle {
	bundle := &DetectionBundle{
		FrameNumber:       frame.GetFrameNumber(),
		CameraID:          frame.GetCameraId(),
		CaptureTime:       frame.GetTCapture(),
		SentTime:          frame.GetTSent(),
		CaptureCameraTime: frame.GetTCaptureCamera(),

		RawBalls:  frame.GetBalls(),
		RawYellow: frame.GetRobotsYellow(),
		RawBlue:   frame.GetRobotsBlue(),

		ReceivedTime: time.Now().UnixMilli(),
		Robots: map[info.Team]map[info.ID]*ssl_vision.SSL_DetectionRobot{
			info.Yellow: {},
			info.Blue:   {},
		},
	}

	for _, r := range bundle.RawYellow {
		id := info.ID(r.GetRobotId())
		bundle.Robots[info.Yellow][id] = r
	}
	for _, r := range bundle.RawBlue {
		id := info.ID(r.GetRobotId())
		bundle.Robots[info.Blue][id] = r
	}

	return bundle
}

