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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/gorilla/websocket"
)

//----------------------------------------------------------------------------------------------
// Start of WebServer class
//----------------------------------------------------------------------------------------------

// Define the WebServer class
type WebServer struct {
	websocketConnections      []*websocket.Conn
	websocketConnectionsMutex sync.Mutex

	websocketupgrader *websocket.Upgrader

	gameStatePacketQueue []([]byte)
	incomingActions      []action.ActionDTO
	gameStateQueueMutex  sync.Mutex
	broadcastThreadMutex sync.Mutex
	receivedDataMutex    sync.Mutex
}

var (
	webserverInstance *WebServer
	Once              sync.Once
)

// Method to get the singleton instance of the WebServer class
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
	go webserverInstance.receiveData()
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

// Method to send the game state to all connected clients
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

// Method to receive data from all connected clients
func (server *WebServer) receiveData() {
	server.websocketConnectionsMutex.Lock()
	defer server.websocketConnectionsMutex.Unlock()

	for _, ws := range server.websocketConnections {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			continue // Use continue here to move to the next iteration
		}

		var receivedData []action.ActionDTO
		err_unmarshal := json.Unmarshal(message, &receivedData)
		if err_unmarshal != nil {
			log.Println("Error unmarshalling message:", err_unmarshal)
			continue // Use continue here to move to the next iteration
		} else {
			server.receivedDataMutex.Lock()
			// Correctly appending the receivedData slice to incomingActions
			server.incomingActions = append(server.incomingActions, receivedData...)
			server.receivedDataMutex.Unlock()
		}
	}
}

//----------------------------------------------------------------------------------------------
// End of WebServer class
//----------------------------------------------------------------------------------------------

// Returns a list of all new incoming actions
func GetIncoming() []action.ActionDTO {
	webserver := GetInstance()
	webserver.receivedDataMutex.Lock()
	defer webserver.receivedDataMutex.Unlock()
	// Return a copy of the incomingActions slice
	actionsCopy := make([]action.ActionDTO, len(webserver.incomingActions))
	copy(actionsCopy, webserver.incomingActions)
	webserver.incomingActions = nil // Empty the incomingActions slice
	return actionsCopy
}

// Broadcasts the game state to all connected clients
func BroadcastGameState(gameStateJson []byte) {
	webserver := GetInstance()
	webserver.gameStateQueueMutex.Lock()
	webserver.gameStatePacketQueue = append(webserver.gameStatePacketQueue, gameStateJson)
	webserver.gameStateQueueMutex.Unlock()
}
