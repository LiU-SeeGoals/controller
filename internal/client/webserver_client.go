package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/info"
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
	// broadcastThreadMutex sync.Mutex
	receivedDataMutex sync.Mutex
}

var (
	webserverInstance *WebServer
	Once              sync.Once
)

// Method to get the singleton instance of the WebServer class
func getInstance() *WebServer {
	Once.Do(startWebServer)
	return webserverInstance
}

// Constructor for the WebServer class
func startWebServer() {
	webserverInstance = &WebServer{
		gameStatePacketQueue: make([]([]byte), 0),
	}

	webserverInstance.websocketupgrader = webserverInstance.getUpgrader()

	http.HandleFunc("/ws", webserverInstance.handleGameStateRequest)
	go http.ListenAndServe(":8080", nil)
	go webserverInstance.sendGameState()
	go webserverInstance.receiveData()
	fmt.Println("server online")
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
	defer server.websocketConnectionsMutex.Unlock() // unlock after function returns
	server.websocketConnections = append(server.websocketConnections, ws)
	fmt.Println("making a connection")
	fmt.Println(len(server.websocketConnections))
	fmt.Print("done serving client")
}

// Method to send the game state to all connected clients
func (server *WebServer) sendGameState() {
	var gameStateJSON []byte
	for {
		if len(server.gameStatePacketQueue) == 0 {
			time.Sleep(time.Millisecond * 10) // Sleep for a short period
			continue
		}

		server.gameStateQueueMutex.Lock()
		gameStateJSON = server.gameStatePacketQueue[0]
		server.gameStatePacketQueue = server.gameStatePacketQueue[1:]
		server.gameStateQueueMutex.Unlock()

		// Creating a copy of the connections. This prevents locking other threads if the connection takes too long
		server.websocketConnectionsMutex.Lock()
		connectionsCopy := make([]*websocket.Conn, len(server.websocketConnections))
		copy(connectionsCopy, server.websocketConnections)
		server.websocketConnectionsMutex.Unlock()

		for _, ws := range connectionsCopy {
			ws.WriteMessage(websocket.TextMessage, gameStateJSON)
			// fmt.Println("written msg")

		}
	}
}

// Method to receive data from all connected clients
func (server *WebServer) receiveData() {
	var validConnections []*websocket.Conn
	for {
		validConnections = validConnections[:0] // reset list

		// Creating a copy of the connections. This prevents locking other threads if the connection takes too long
		server.websocketConnectionsMutex.Lock()
		connectionsCopy := make([]*websocket.Conn, len(server.websocketConnections))
		copy(connectionsCopy, server.websocketConnections)
		server.websocketConnectionsMutex.Unlock()

		for _, ws := range connectionsCopy {
			_, message, err := ws.ReadMessage()
			if err != nil {
				ws.Close()
				continue
			}

			var receivedData action.ActionDTO
			err_unmarshal := json.Unmarshal(message, &receivedData)
			if err_unmarshal != nil {
				log.Println("Error unmarshalling message:", err_unmarshal)
				continue
			} else {
				server.receivedDataMutex.Lock()
				log.Println("Received data:", receivedData)
				server.incomingActions = append(server.incomingActions, receivedData)
				server.receivedDataMutex.Unlock()
			}
			validConnections = append(validConnections, ws)
		}

		server.websocketConnectionsMutex.Lock()
		// Remove invalid connections
		server.websocketConnections = validConnections
		server.websocketConnectionsMutex.Unlock()
	}
}

//----------------------------------------------------------------------------------------------
// End of WebServer class
//----------------------------------------------------------------------------------------------

// How to use the WebServer class:
// Only use the functions under this comment to interact with the WebServer class
// The WebServer class is a singleton class, so you can only have one instance of it,
// and the functions under handles all of it so multiple instances are not created

type WebsiteDTO struct {
	RobotPositions [2 * info.TEAM_SIZE]info.RobotDTO
	BallPosition   info.BallDTO
	RobotActions   []action.ActionDTO
	TerminalLog    []string
}

func toJson(input WebsiteDTO) []byte {
	output, err := json.Marshal(input)
	if err != nil {
		fmt.Println("The WebsiteDTO packet could not be marshalled to JSON.")
	}
	return output
}

// Returns a list of all new incoming actions
func GetIncoming() []action.ActionDTO {
	webserver := getInstance()
	webserver.receivedDataMutex.Lock()
	defer webserver.receivedDataMutex.Unlock()
	// Return a copy of the incomingActions slice
	actionsCopy := make([]action.ActionDTO, len(webserver.incomingActions))
	copy(actionsCopy, webserver.incomingActions)
	webserver.incomingActions = nil // Empty the incomingActions slice
	return actionsCopy
}

// Broadcasts the game state to all connected clients
func BroadcastGameState(message WebsiteDTO) {
	gameStateJson := toJson(message)
	webserver := getInstance()
	webserver.gameStateQueueMutex.Lock()
	webserver.gameStatePacketQueue = append(webserver.gameStatePacketQueue, gameStateJson)
	webserver.gameStateQueueMutex.Unlock()
}

func UpdateWebGUI(gs *info.GameState, actions []action.Action, terminal_messages []string) {
	var gamestate_DTO = gs.ToDTO()
	var actionTDO = make([]action.ActionDTO, len(actions))
	for i, obj := range actions {
		actionTDO[i] = obj.ToDTO()
	}
	var websiteMessage = WebsiteDTO{
		RobotPositions: gamestate_DTO.RobotPositions,
		BallPosition:   gamestate_DTO.BallPosition,
		RobotActions:   actionTDO,
		TerminalLog:    terminal_messages,
	}
	fmt.Println(websiteMessage.RobotActions)
	BroadcastGameState(websiteMessage)
}
