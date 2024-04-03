package client

import (
	"log/slog"
	"net"
	"sync"

	"google.golang.org/protobuf/proto"
)

// `BaseStationClient` is a UDP client that sends protobuf messages of a generic type.
// It maintains a connection to a base station and ensures thread-safe access to this connection.
type BaseStationClient[T proto.Message] struct {
	connection    net.Conn
	connMutex	  sync.Mutex
}

// `NewBaseStationClient`` creates a new BaseStationClient with a connection to the given address.
// It establishes a UDP connection to the base station and returns the client.
//
// Parameters:
// 	address: The address of the base station to connect to.
//
// Returns:
// 	A pointer to the created BaseStationClient or nil if the connection failed.
func NewBaseStationClient[T proto.Message](address string) *BaseStationClient[T] {
	slog.Info("Creating new base station client on", slog.String("address", address))
	connection, err := net.Dial("udp", address)
	if err != nil {
		slog.Error("Failed to dial UDP connection on address %s.\n%v", address, err)
		return nil
	}

	return &BaseStationClient[T]{
		connection: connection,
		connMutex: sync.Mutex{},
	}
}

// The fucntion `Send` asynchronously sends a protobuf message to the base station.
// This function is non-blocking due to calling the `sendItem` function in a new goroutine.
func (b *BaseStationClient[T]) Send(item T) {
	go b.sendItem(item)
}

// `sendItem` serializes the given protobuf message and sends it over the connection.
// It ensures thread-safe access to the connection. 
// If serialization or sending fails, it logs an error. It also logs if not all bytes were sent.
func (b *BaseStationClient[T]) sendItem(item T) {
	var serializedItem []byte
	var err error
	var n_bytes_sent int

	serializedItem, err = proto.Marshal(item)
	if err != nil {
		slog.Error("Failed to serialize item. ", slog.String("err", err.Error()))
		return
	}

	b.connMutex.Lock()
	n_bytes_sent, err = b.connection.Write(serializedItem)
	b.connMutex.Unlock()

	if err != nil {
		slog.Warn("Failed to send item. ", slog.String("err", err.Error()))
	}
	if n_bytes_sent != len(serializedItem) {
		slog.Warn("Failed to send all bytes of item.")
	}
}