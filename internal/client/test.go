package client

// package main
//
// import (
// 	"fmt"
// 	"log"
//
// 	"github.com/LiU-SeeGoals/controller/internal/proto/basestation" // This import path assumes 'controller' is the module name
//
// 	"google.golang.org/protobuf/proto"
// )
//
// func main() {
// 	// Create a new Command instance
// 	command_move := &basestation.Command{
// 		CommandId: basestation.ActionType_MOVE_ACTION, // Setting the command type
// 		RobotId:   1,                                  // Assuming a robot ID
// 		Pos: &basestation.Vector3D{                    // Set the position
// 			X: 1,
// 			Y: 1,
// 			W: 30.5,
// 		},
// 		Goal: &basestation.Vector3D{                   // Set the goal
// 			X: 1,
// 			Y: 1,
// 			W: 70.5,
// 		},
// 	}
//
// 	command_stop := &basestation.Command{
// 		CommandId: basestation.ActionType_STOP_ACTION, // Setting the command type
// 		RobotId:   2,                                  // Assuming a robot ID
// 	}
//
// 	// Serialize the Command to a byte slice
// 	data, err := proto.Marshal(command_move)
// 	if err != nil {
// 		log.Fatalf("Failed to marshal command: %v", err)
// 	}
//
// 	fmt.Printf("Serialized data move: %x\n", data)
//
// 	// Serialize the Command to a byte slice
// 	data, err = proto.Marshal(command_stop)
// 	if err != nil {
// 		log.Fatalf("Failed to marshal command: %v", err)
// 	}
//
// 	fmt.Printf("Serialized data stop: %x\n", data)
// }