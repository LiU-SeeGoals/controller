package logger

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger  *zap.Logger
	LoggerS *zap.SugaredLogger
)

// Implements io.Writer
type WebWriter struct{
	server *logServer
}

func (w *WebWriter) Write(p []byte) (n int, err error) {
	UpdateWebLog(p, w.server)
	return len(p), nil
}

func init() {

	webserverInstance := &logServer{
		logPacketQueue: make([]([]byte), 0),
	}
	startWebServer(webserverInstance)

	// AddSync converts an io.Writer to a WriteSyncer.
	writeSyncer := zapcore.AddSync(&WebWriter{server: webserverInstance})

	// Lock wraps a WriteSyncer in a mutex to make it safe for concurrent use.
	// In particular, *os.Files must be locked before use.
	writeSyncer = zapcore.Lock(writeSyncer)

	// NewMultiWriteSyncer creates a WriteSyncer that duplicates its writes
	// and sync calls, much like io.MultiWriter.
	multiWriter := zapcore.NewMultiWriteSyncer(writeSyncer, os.Stdout)

	lvl := zapcore.DebugLevel
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, multiWriter, lvl)
	Logger = zap.New(core, zap.AddCaller())

	// Sugar wraps the Logger to provide a more ergonomic, but slightly slower, API.
	LoggerS = Logger.Sugar()
}

//----------------------------------------------------------------------------------------------
// Start of logServer class
//----------------------------------------------------------------------------------------------

// Define the logServer class
type logServer struct {
	websocketConnections      []*websocket.Conn
	websocketConnectionsMutex sync.Mutex

	websocketupgrader *websocket.Upgrader

	logPacketQueue []([]byte)
	logQueueMutex  sync.Mutex
}

// Constructor for the logServer class
func startWebServer(webserverInstance *logServer) {

	webserverInstance.websocketupgrader = webserverInstance.getUpgrader()

	http.HandleFunc("/logs", webserverInstance.handleGameStateRequest)
	go http.ListenAndServe(":8080", nil)
	go webserverInstance.sendLog()
	fmt.Println("log server online")
}

func (server *logServer) getUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

// Method to handle connections
func (server *logServer) handleGameStateRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("serving client")
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

// Method to send the logs to all connected clients
func (server *logServer) sendLog() {
	var logJSON []byte
	for {
		if len(server.logPacketQueue) == 0 {
			time.Sleep(time.Millisecond * 10) // Sleep for a short period
			continue
		}

		server.logQueueMutex.Lock()
		logJSON = server.logPacketQueue[0]
		server.logPacketQueue = server.logPacketQueue[1:]
		server.logQueueMutex.Unlock()

		// Creating a copy of the connections. This prevents locking other threads if the connection takes too long
		server.websocketConnectionsMutex.Lock()
		connectionsCopy := make([]*websocket.Conn, len(server.websocketConnections))
		copy(connectionsCopy, server.websocketConnections)
		server.websocketConnectionsMutex.Unlock()

		for _, ws := range connectionsCopy {
			ws.WriteMessage(websocket.TextMessage, logJSON)
			// fmt.Println("written msg")

		}
	}
}

// Sends logs to the webGUI
func UpdateWebLog(log []byte, webserver *logServer) {

	// Lock the mutex to protect access to the log queue
	webserver.logQueueMutex.Lock()
	defer webserver.logQueueMutex.Unlock()

	// Add the log entry to the log queue
	webserver.logPacketQueue = append(webserver.logPacketQueue, log)
}

