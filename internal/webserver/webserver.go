// package webserver
//
// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
//
// 	"github.com/LiU-SeeGoals/controller/internal/gamestate"
// 	"github.com/gorilla/websocket"
// )
//
//
// func HandleConnections(w http.ResponseWriter, r *http.Request) {
// 	upgrader := websocket.Upgrader{
// 		CheckOrigin: func(r *http.Request) bool {
// 			return true // Allow connections from any origin
// 		},
// 	}
//
//     // Upgrade initial GET request to a WebSocket
//     ws, err := upgrader.Upgrade(w, r, nil)
//     if err != nil {
//         panic(err)
//     }
//     defer ws.Close()
//
//     // Initialize your game state here or pass it to this function
//     gs := gamestate.NewGameState("127.0.0.1:20011", "224.5.23.2:10020")
//
//     // Main loop: wait for a message and send it to all connected clients
//     for {
//         // Read message from browser (you can also remove this if you don't expect messages from the client)
//
//         // Update the game state
//         gs.Update() // Assuming this updates your game state
//
//         gameStateJSON, err := json.Marshal(gs.ToDTO())
//         // Serialize the game state to JSON
//         if err != nil {
//             panic(err) // Handle the error appropriately
//         }
// 		fmt.Println(string(gameStateJSON))
// 		ws.WriteMessage(websocket.TextMessage, gameStateJSON);
//
//         // Send the serialized game state to the client
//         if err := ws.WriteMessage(websocket.TextMessage, gameStateJSON); err != nil {
//             break // Exit the loop if there's an error
//         }
//     }
// }

package webserver

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Define the WebServer class
type WebServer struct {
	websocketConnections      []*websocket.Conn
	websocketConnectionsMutex sync.Mutex

	websocketupgrader *websocket.Upgrader

	gameStatePacketQueue []([]byte)
	gameStateQueueMutex  sync.Mutex
	broadcastThreadMutex sync.Mutex
}

var (
	webserverInstance *WebServer
	Once              sync.Once
)

func GetInstance() *WebServer {
	Once.Do(StartWebServer)
	return webserverInstance
}

// Constructor for the WebServer class
func StartWebServer() {
	webserverInstance = &WebServer{
		gameStatePacketQueue: make([]([]byte), 0),
	}

	webserverInstance.websocketupgrader = webserverInstance.getUpgrader()

	http.HandleFunc("/ws", webserverInstance.handleGameStateRequest)
	go http.ListenAndServe(":8080", nil)
	go webserverInstance.sendGameState()
}

func BroadcastGameState(gameStateJson []byte) {
	webserver := GetInstance()
	webserver.gameStateQueueMutex.Lock()
	webserver.gameStatePacketQueue = append(webserver.gameStatePacketQueue, gameStateJson)
	webserver.gameStateQueueMutex.Unlock()
}

func (server *WebServer) getUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

// Method to handle connections
func (server *WebServer) handleGameStateRequest(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a WebSocket
	ws, err := server.websocketupgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	server.websocketConnectionsMutex.Lock()
	server.websocketConnections = append(server.websocketConnections, ws)
	fmt.Println("making a connection")
	fmt.Println(len(server.websocketConnections))
	b := make([]byte, 10)
	ws.WriteMessage(websocket.TextMessage, b)
	server.websocketConnectionsMutex.Unlock()
	fmt.Print("done serving client")
}

var (
	helloman int
)

func (server *WebServer) sendGameState() {
	var gameStateJSON []byte
	for {
		if len(server.gameStatePacketQueue) == 0 {
			continue
		}

		server.gameStateQueueMutex.Lock()
		gameStateJSON = server.gameStatePacketQueue[0]
		server.gameStatePacketQueue = server.gameStatePacketQueue[1:]
		server.gameStateQueueMutex.Unlock()

		server.websocketConnectionsMutex.Lock()
		for _, ws := range server.websocketConnections {
			fmt.Println("writing")
			//fmt.Println(str(gameStateJSON))
			ws.WriteMessage(websocket.TextMessage, gameStateJSON)
		}
		server.websocketConnectionsMutex.Unlock()
	}
}
