package receiver

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/LiU-SeeGoals/proto_go/ssl_vision"
	"github.com/golang/protobuf/proto"
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
	// Channel to send packets to
	ssl_receiver_channel chan *ssl_vision.SSL_WrapperPacket
}

// Create a new SSL vision receiver.
// Address should be <ip>:<port> - grSim default is 224.5.23.2:10020
func NewSSLReceiver(addr string) *SSLReceiver {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		slog.Error("Unable to resolve UDP address: %s", err)
		panic(err)
	}

	return &SSLReceiver{
		conn: nil,
		addr: udpAddr,
		buff: make([]byte, read_buff_sz),
		ssl_receiver_channel: make(chan *ssl_vision.SSL_WrapperPacket),
	}
}

// Connect/subscribe receiver to UDP multicast.
// Note, this will NOT block.
func (r *SSLReceiver) Connect() {
	conn, err := net.ListenMulticastUDP("udp", nil, r.addr)
	if err != nil {
		slog.Error("Unable to connect to UDP multicast: %s", err)
		panic(err)
	}

	r.conn = conn
	go r.receivePackets()
}

func (r *SSLReceiver) Receive() *ssl_vision.SSL_WrapperPacket {
	packet := <-r.ssl_receiver_channel
	return packet
}

// Start receiving packets.
// This function should be run in a goroutine:
//
//	go recv.Receive()
//
// Parsed packets are transferred using packetChan.
func (r *SSLReceiver) receivePackets() {
	var packet *ssl_vision.SSL_WrapperPacket
	for {
		packet = &ssl_vision.SSL_WrapperPacket{}
		sz, err := r.conn.Read(r.buff)
		if err != nil {
			fmt.Printf("Unable to receive packet: %s", err)
			continue
		}

		err = proto.Unmarshal(r.buff[:sz], packet)
		if err != nil {
			fmt.Printf("Unable to unmarshal packet: %s", err)
			continue
		}

		r.ssl_receiver_channel <- packet
	}
}
