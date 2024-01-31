package main

import (
	"fmt"
	"net"

	"github.com/LiU-SeeGoals/controller/internal/parsed_vision"

	"google.golang.org/protobuf/proto"
)

func main() {
	// Just change this two variable
	// There is no error handling
	// Good luck :)
	parsedFrame := parsed_vision.ParsedFrame{} 
	address := "127.0.0.1:1234"

	var err error = nil
	connection, _ := net.Dial("udp", address)
	if err != nil {
		fmt.Printf("Some error %v\n", err)
	}

	bytes, _ := proto.Marshal(&parsedFrame)
	connection.Write(bytes)
}
