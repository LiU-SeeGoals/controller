package client

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/proto/basestation"
	"github.com/golang/protobuf/proto"
)

const MAX_SEND_SIZE = 2048

type Connection interface {
    Write(b []byte) (n int, err error)
    Close() error
}

type BaseStationClient struct {
	connection Connection
	address string
	queueMutex sync.Mutex
	threadMutex sync.Mutex
	queue []*basestation.Command
	hasBeenInited bool

}

func NewBaseStationClient(address string) *BaseStationClient {
	var err error = nil
	connection, _ := net.Dial("udp", address)
	if err != nil {
        fmt.Printf("Some error %v\n", err)
    }

	return &BaseStationClient{
		connection: connection,
		address: address,
		queue: make([]*basestation.Command, 0),
		hasBeenInited: false,
	}	
}

func (b *BaseStationClient) Init() {
	fmt.Printf("running init\n")
	go b.sendCommands()
	b.hasBeenInited = true
	b.threadMutex.Lock()
}

// Goroutine function for processing and sending commands
func (b *BaseStationClient) sendCommands() {
	fmt.Printf("starting up the send command\n")
    for {

        if len(b.queue) == 0 {
            // Wait to be unlocked
            b.threadMutex.Lock()
			continue
        }
	
        // Process the first command in the queue
		b.queueMutex.Lock()
        cmd := b.queue[0]
		b.queue = b.queue[1:]
		b.queueMutex.Unlock()
		
		fmt.Printf("Sending the shit\n")
		
        // Send the command
        serializedCmd, _ := proto.Marshal(cmd) // Add error handling
        b.sendMessage(serializedCmd)           // Add error handling
    }
}


func (b *BaseStationClient) Send(actions []action.Action) {
	if !b.hasBeenInited{
		fmt.Println( "\033[0m Base station client has not been inited\033[33m")
	}

	for _, action := range actions {
		fmt.Printf("Adding the shit to the queue.\n")
		b.queueMutex.Lock()
		b.queue = append(b.queue, action.TranslateReal())
		b.queueMutex.Unlock()
	}
	b.threadMutex.Unlock()
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