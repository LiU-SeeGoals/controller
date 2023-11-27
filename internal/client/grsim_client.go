package client

import (
	"fmt"
	"net"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/datatypes"
	"github.com/LiU-SeeGoals/controller/internal/proto/grsim"
	"github.com/golang/protobuf/proto"
)

// SSL Vision receiver
type GrsimClient struct {
	// Connection
	conn *net.UDPConn

	// UDP address
	addr *net.UDPAddr

	// Yellow team robot command buffer
	buffYellow []*grsim.GrSim_Robot_Command

	// Blue team robot command buffer
	buffBlue []*grsim.GrSim_Robot_Command

	// Local time
	// Note: grsim requires us to send a "timestamp",
	// but it's really unclear what this timestamp is.
	// For now, it's implemented as a local counter.
	time float64
}

// Create new Grsim client
// Address should be <ip>:<port>
func NewSSLGrsimClient(addr string) *GrsimClient {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}

	return &GrsimClient{
		conn: nil,
		addr: udpAddr,
	}
}

// Connect/subscribe receiver to UDP multicast.
// Note, this will NOT block.
func (client *GrsimClient) Connect() {
	conn, err := net.DialUDP("udp", nil, client.addr)
	if err != nil {
		panic(err)
	}

	client.conn = conn
}

func (client *GrsimClient) AddActions(actions []action.Action) {

	for id, action := range actions {
		params := datatypes.NewParameters()
		params.RobotId = uint32(id)
		action.TranslateGrsim(params)
		client.addRobotCommand(params)
	}
	client.send()
}

// Add a new Robot command to client buffer
func (client *GrsimClient) addRobotCommand(params *datatypes.Parameters) {
	command := newRobotCommand(params.RobotId, params.VelTangent, params.VelNormal, params.VelAngular, params.KickSpeedX, params.KickSpeedZ, params.Spinner, params.WheelsSpeed)

	if params.YellowTeam {
		client.buffYellow = append(client.buffYellow, command)
		return
	}
	client.buffBlue = append(client.buffBlue, command)
}

// Helper function creates a new robot command
func newRobotCommand(
	robotId uint32,
	velTangent float32,
	velNormal float32,
	velAngular float32,
	kickSpeedX float32,
	kickSpeedZ float32,
	spinner bool,
	wheelsSpeed bool,
) *grsim.GrSim_Robot_Command {
	return &grsim.GrSim_Robot_Command{
		Id:         &robotId,
		Kickspeedx: &kickSpeedX,
		Kickspeedz: &kickSpeedZ,

		Veltangent: &velTangent,
		Velnormal:  &velNormal,
		Velangular: &velAngular,

		Spinner:     &spinner,
		Wheelsspeed: &wheelsSpeed,
	}
}

// Helper function clears command buffer
func (client *GrsimClient) clearCommandBuffer() {
	client.buffYellow = []*grsim.GrSim_Robot_Command{}
	client.buffBlue = []*grsim.GrSim_Robot_Command{}
}

func (client *GrsimClient) send() (int, error) {
	// Incr time
	client.time += 1.0

	// Clear buffers
	defer client.clearCommandBuffer()

	packet := &grsim.GrSim_Packet{}

	isteamyellow := true
	packet.Commands = &grsim.GrSim_Commands{
		Timestamp:     &client.time,
		Isteamyellow:  &isteamyellow,
		RobotCommands: client.buffYellow,
	}

	// Yellow team
	data, err := proto.Marshal(packet)
	if err != nil {
		err = fmt.Errorf("unable to marshal yellow team data: %w", err)
		return 0, err
	}

	writeYellow, err := client.conn.Write(data)
	if err != nil {
		err = fmt.Errorf("unable to send yellow team data over socket: %w", err)
		return 0, err
	}

	isteamyellow = false
	packet.Commands = &grsim.GrSim_Commands{
		Timestamp:     &client.time,
		Isteamyellow:  &isteamyellow,
		RobotCommands: client.buffBlue,
	}

	// Blue team
	data, err = proto.Marshal(packet)
	if err != nil {
		err = fmt.Errorf("unable to marshal blue team data: %w", err)
		return 0, err
	}

	writeBlue, err := client.conn.Write(data)
	if err != nil {
		err = fmt.Errorf("unable to send blue team data over socket: %w", err)
		return 0, err
	}

	return writeYellow + writeBlue, nil
}
