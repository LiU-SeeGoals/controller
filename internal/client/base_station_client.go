package client

import (
	"fmt"
	"log/slog"
	"net"
	"sync"

	"google.golang.org/protobuf/proto"
)

type BaseStationClient[T proto.Message] struct {
	connection    net.Conn
	connMutex	  sync.Mutex
}

func NewBaseStationClient[T proto.Message](address string, port int) *BaseStationClient[T] {

	fullAddress := fmt.Sprintf("%s:%d", address, port)
	connection, err := net.Dial("udp", fullAddress)
	if err != nil {
		slog.Error("Failed to dial UDP connection on address %s.\n%v", address, err)
		return nil
	}

	return &BaseStationClient[T]{
		connection: connection,
		connMutex: sync.Mutex{},
	}
}

func (b *BaseStationClient[T]) Send(item T) {
	go b.sendItem(item)
}

func (b *BaseStationClient[T]) sendItem(item T) {
	serializedItem, err := proto.Marshal(item)
	fmt.Println(serializedItem)
	if err != nil {
		slog.Error("Failed to serialize item.\n%v", err)
		return
	}

	var n_bytes_sent int
	b.connMutex.Lock()
	n_bytes_sent, err = b.connection.Write(serializedItem)
	b.connMutex.Unlock()

	if err != nil {
		slog.Error("Failed to send item.\n%v", err)
	}
	if n_bytes_sent != len(serializedItem) {
		slog.Error("Failed to send all bytes of item.\n")
	}
}