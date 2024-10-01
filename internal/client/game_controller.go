package receiver

import (
	"fmt"
	"net"

	"github.com/LiU-SeeGoals/proto_go/gc"
	"google.golang.org/protobuf/proto"
)

// const (
// 	// Read buffer size
// 	READ_BUFFER_SIZE = 8192
// )

// SSL Vision receiver
type GCConnection struct {
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
func NewGCConnection(addr string) *GCConnection {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}

	return &GCConnection{
		conn: nil,
		addr: udpAddr,
		buff: make([]byte, READ_BUFFER_SIZE),
	}
}

// Connect/subscribe receiver to UDP multicast.
// Note, this will NOT block.
func (r *GCConnection) Connect() {
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
func (r *GCConnection) Receive(packetChan chan *gc.Referee) {
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

type GCReceiver struct {
	gc         *GCConnection
	gc_channel chan *gc.Referee
}

// Start a SSL Vision receiver, returns a channel from
// which SSL wrapper packets can be obtained.
func (receiver *GCReceiver) Connect() {
	receiver.gc.Connect()
	go receiver.gc.Receive(receiver.gc_channel)
}

func NewGCReceiver(sslReceiverAddress string) *GCReceiver {
	receiver := &GCReceiver{
		gc:         NewGCConnection(sslReceiverAddress),
		gc_channel: make(chan *gc.Referee),
	}
	receiver.Connect()
	return receiver
}

// Test printing out packets
func (receiver *GCReceiver) PrintPackets() {
	for {
		packet, ok := <-receiver.gc_channel
		if !ok {
			fmt.Println("GC Channel closed")
			return
		}
		fmt.Println(packet)
	}
}

