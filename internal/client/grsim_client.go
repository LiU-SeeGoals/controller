package client

import (
	"net"
	"fmt"
)


// SSL Vision receiver
type GrsimClient struct {
	// Connection
	conn *net.UDPConn
	// UDP address
	addr *net.UDPAddr
}

func NewSSLGrsimClient(addr string) GrsimClient {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}

	return GrsimClient{
		conn: nil,
		addr: udpAddr,
	}
}

func (client *GrsimClient) Connect() {
	conn, err := net.DialUDP("udp", nil, client.addr)
	if err != nil {
		panic(err)
	}

	client.conn = conn
}

func (client *GrsimClient) Send()  {
	writtenBytes, err := client.conn.Write([]byte{123,123})

	if err != nil {
		fmt.Printf("It's a fucking error")
		return
	}

	fmt.Println("It fucking worked")
	fmt.Println(writtenBytes)
}

