package client

import (
	"errors"
	"fmt"
	"net"

	"github.com/LiU-SeeGoals/controller/internal/action"
)

const MAX_SEND_SIZE = 2048

type Connection interface {
    Write(b []byte) (n int, err error)
    Close() error
}

type BaseStationClient struct {
	connection Connection
	address string
	
}

func NewBaseStationClient(address string) *BaseStationClient {
	connection, _ := net.Dial("udp", address)
	return &BaseStationClient{
		connection: connection,
		address: address,
	}	
}

func (b *BaseStationClient) OpenConnection() error{
	var err error = nil
	b.connection, err = net.Dial("udp", b.address)
	if err != nil {
        fmt.Printf("Some error %v\n", err)
        return err
    }
	return err
}

func (b *BaseStationClient) SendActions(actions []action.Action) {
	// The structure of a message:
	// 1. The length of the message in bytes
	// 2. The action type for the robot
	// 3. The robot to perform the action
	// 4. Params
	for _, action := range actions {
		b.sendMessage(action.TranslateReal())
	}
}

func (b *BaseStationClient) sendMessage(input []byte) error{
	if (len(input) > MAX_SEND_SIZE) {
		fmt.Print("to big to send (if sent = Rasmus mad 😡)")
		return errors.New("too long message")
	}

	_,err := b.connection.Write(input)
    if err != nil {
        fmt.Printf("Some error %v\n", err)
		return err
    }
	return nil
}

func (b *BaseStationClient) CloseConnection() {
	b.connection.Close()	
}