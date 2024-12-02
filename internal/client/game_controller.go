package client

import (
	"fmt"
	"net"

	"github.com/LiU-SeeGoals/controller/internal/gamestatus"
	"github.com/LiU-SeeGoals/proto_go/gc"
	"google.golang.org/protobuf/proto"
)

// SSL Vision receiver
type GCConnection struct {
	// Connection
	conn *net.UDPConn
	// UDP address
	addr *net.UDPAddr
	// Read buffer
	buff []byte
	// SSL lets not heap allocate this every time
	packet gc.Referee
}

// Create a new SSL vision receiver.
// Address should be <ip>:<port>
func NewGCConnection(addr string) *GCConnection {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}

	return &GCConnection{
		conn: nil,
		addr: udpAddr,
		buff: make([]byte, READ_BUFFER_SIZE),
	}
}

// Connect/subscribe receiver to UDP multicast.
// Note, this will NOT block.
func (r *GCConnection) Connect() {
	conn, err := net.ListenMulticastUDP("udp", nil, r.addr)
	if err != nil {
		panic(err)
	}

	r.conn = conn
}

// Start receiving packets.
// This function should be run in a goroutine:
//
//	go recv.Receive()
//
// Parsed packets are transferred using packetChan.
func (r *GCConnection) Receive(packetChan chan *gc.Referee) {
	for {
		sz, err := r.conn.Read(r.buff)
		if err != nil {
			fmt.Printf("Unable to receive packet: %s", err)
			continue
		}

		err = proto.Unmarshal(r.buff[:sz], &r.packet)
		if err != nil {
			fmt.Printf("Unable to unmarshal packet: %s", err)
			continue
		}

		packetChan <- &r.packet
	}
}

type GCClient struct {
	gc         *GCConnection
	gc_channel chan *gc.Referee
	// A random UUID of the source that is kept constant at the source while running
	// If multiple sources are broadcasting to the same network, this id can be used to identify individual sources
	SourceIdentifier string
}

// Start a SSL Vision receiver, returns a channel from
// which SSL wrapper packets can be obtained.
func (receiver *GCClient) Connect() {
	receiver.gc.Connect()
	go receiver.gc.Receive(receiver.gc_channel)
}

func NewGCClient(sslReceiverAddress string) *GCClient {
	receiver := &GCClient{
		gc:         NewGCConnection(sslReceiverAddress),
		gc_channel: make(chan *gc.Referee),
	}
	receiver.Connect()
	return receiver
}

func (receiver *GCClient) InitGameStatus(gs *gamestatus.GameStatus) {
	packet, ok := <-receiver.gc_channel

	if !ok {
		fmt.Println("GC Channel closed")
		return
	}
	receiver.SourceIdentifier = packet.GetSourceIdentifier()
}

// Test printing out packets
func (receiver *GCClient) UpdateGameStatus(gs *gamestatus.GameStatus) {
	packet, ok := <-receiver.gc_channel

	if packet.GetSourceIdentifier() != receiver.SourceIdentifier {
		return
	}

	if !ok {
		fmt.Println("GC Channel closed")
		return
	}

	gs.SetGameEvent(gamestatus.RefCommand(packet.GetCommand().Number()),
		packet.GetCommandTimestamp(),
		float64(packet.GetDesignatedPosition().GetX()),
		float64(packet.GetDesignatedPosition().GetY()),
		gamestatus.RefCommand(packet.GetCommand().Number()),
		packet.GetCurrentActionTimeRemaining())

	gs.SetGameStatus(gamestatus.GameStage(packet.GetStage().Number()),
		gamestatus.MatchType(packet.GetMatchType().Number()),
		packet.GetPacketTimestamp(),
		packet.GetStageTimeLeft(),
		packet.GetCommandCounter(),
		packet.GetBlueTeamOnPositiveHalf(),
		packet.GetStatusMessage())

	// yellow team
	gs.SetTeamInfo(
		true,
		packet.Yellow.GetName(),
		packet.Yellow.GetScore(),
		packet.Yellow.GetRedCards(),
		packet.Yellow.GetYellowCards(),
		packet.Yellow.GetTimeouts(),
		packet.Yellow.GetTimeoutTime(),
		packet.Yellow.GetGoalkeeper(),
		packet.Yellow.GetFoulCounter(),
		packet.Yellow.GetBallPlacementFailures(),
		packet.Yellow.GetMaxAllowedBots(),
		packet.Yellow.GetBotSubstitutionsLeft(),
		packet.Yellow.GetBotSubstitutionTimeLeft(),
		packet.Yellow.GetYellowCardTimes(),
		packet.Yellow.GetCanPlaceBall(),
		packet.Yellow.GetBotSubstitutionIntent(),
		packet.Yellow.GetBallPlacementFailuresReached(),
		packet.Yellow.GetBotSubstitutionAllowed(),
	)

	// blue team
	gs.SetTeamInfo(
		false,
		packet.Blue.GetName(),
		packet.Blue.GetScore(),
		packet.Blue.GetRedCards(),
		packet.Blue.GetYellowCards(),
		packet.Blue.GetTimeouts(),
		packet.Blue.GetTimeoutTime(),
		packet.Blue.GetGoalkeeper(),
		packet.Blue.GetFoulCounter(),
		packet.Blue.GetBallPlacementFailures(),
		packet.Blue.GetMaxAllowedBots(),
		packet.Blue.GetBotSubstitutionsLeft(),
		packet.Blue.GetBotSubstitutionTimeLeft(),
		packet.Blue.GetYellowCardTimes(),
		packet.Blue.GetCanPlaceBall(),
		packet.Blue.GetBotSubstitutionIntent(),
		packet.Blue.GetBallPlacementFailuresReached(),
		packet.Blue.GetBotSubstitutionAllowed(),
	)

}
