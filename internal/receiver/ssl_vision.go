package receiver

import (
	"fmt"
	"net"

	"github.com/LiU-SeeGoals/proto-messages/ssl_vision"
	"google.golang.org/protobuf/proto"
)

const (
	// Read buffer size
	read_buff_sz = 8192
)

// SSL Vision receiver
type SSLReceiver struct {
	// Connection
	conn *net.UDPConn
	// UDP address
	addr *net.UDPAddr
	// Read buffer
	buff []byte
}

// Create a new SSL vision receiver.
// Address should be <ip>:<port> - grSim default is 224.5.23.2:10020
func NewSSLReceiver(addr string) *SSLReceiver {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}

	return &SSLReceiver{
		conn: nil,
		addr: udpAddr,
		buff: make([]byte, read_buff_sz),
	}
}

// Connect/subscribe receiver to UDP multicast.
// Note, this will NOT block.
func (r *SSLReceiver) Connect() {
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
func (r *SSLReceiver) Receive(packetChan chan ssl_vision.SSL_WrapperPacket) {
	var packet ssl_vision.SSL_WrapperPacket
	for {
		sz, err := r.conn.Read(r.buff)
		if err != nil {
			fmt.Printf("Unable to receive packet: %s", err)
			continue
		}

		err = proto.Unmarshal(r.buff[:sz], &packet)
		if err != nil {
			fmt.Printf("Unable to unmarshal packet: %s", err)
			continue
		}

		packetChan <- packet
	}
}
