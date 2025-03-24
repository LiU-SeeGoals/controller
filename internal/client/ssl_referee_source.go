package client

import (
	"fmt"
	"net"

	"github.com/LiU-SeeGoals/controller/internal/helper"
	"github.com/LiU-SeeGoals/proto_go/gc"
	"google.golang.org/protobuf/proto"
)

// SSL Vision receiver
type refereeConnection struct {
	// Connection
	conn *net.UDPConn
	// UDP address
	addr *net.UDPAddr
	// Read buffer
	buff []byte
	// SSL lets not heap allocate this every time
	packet gc.Referee
}

// Create a new SSL vision receiver.
// Address should be <ip>:<port>
func NewRefereeConnection(addr string) *refereeConnection {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}

	return &refereeConnection{
		conn: nil,
		addr: udpAddr,
		buff: make([]byte, READ_BUFFER_SIZE),
	}
}

// Connect/subscribe receiver to UDP multicast.
func (r *refereeConnection) Connect() {
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
func (r *refereeConnection) Receive(packetChan chan *gc.Referee) {
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

		helper.NB_Send[gc.Referee](packetChan, &r.packet)

	}
}

type SSLRefereeSource struct {
	gc         *refereeConnection
	gc_channel chan *gc.Referee
	// A random UUID of the source that is kept constant at the source while running
	// If multiple sources are broadcasting to the same network, this id can be used to identify individual sources
	SourceIdentifier string
}

// Start a SSL Vision receiver, returns a channel from
// which SSL wrapper packets can be obtained.
func (receiver *SSLRefereeSource) Connect() {
	receiver.gc.Connect()
	go receiver.gc.Receive(receiver.gc_channel)
}


func NewSSLRefereeSource(sslReceiverAddress string) *SSLRefereeSource {
	receiver := &SSLRefereeSource{
		gc:         NewRefereeConnection(sslReceiverAddress),
		gc_channel: make(chan *gc.Referee),
	}
	receiver.Connect()
	return receiver
}


func (receiver *SSLRefereeSource) GetRefereeData() *gc.Referee{
	select {
	case packet, ok := <- receiver.gc_channel:
		if !ok {
			fmt.Println("GC Channel closed")
			return nil
		}
		return packet
	default:
		return nil
	}
}

