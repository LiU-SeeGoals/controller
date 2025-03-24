package client

import (
	"net"
	"sync/atomic"

	"github.com/LiU-SeeGoals/controller/internal/info"
	. "github.com/LiU-SeeGoals/controller/internal/logger"
	"github.com/LiU-SeeGoals/proto_go/gc"
	"github.com/LiU-SeeGoals/proto_go/ssl_vision"
	"google.golang.org/protobuf/proto"
)

type AtomicLatest[T any] struct {
	ptr atomic.Pointer[T]
}

func NewAtomicLatest[T any]() *AtomicLatest[T] {
	return &AtomicLatest[T]{}
}

func (a *AtomicLatest[T]) Set(v *T) { a.ptr.Store(v) }
func (a *AtomicLatest[T]) Get() *T  { return a.ptr.Load() }

type trackerConnection struct {
	conn *net.UDPConn
	addr *net.UDPAddr
	buff []byte
}

func NewTrackerConnection(addr string) *trackerConnection {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		panic(err)
	}
	return &trackerConnection{
		conn: nil,
		addr: udpAddr,
		buff: make([]byte, READ_BUFFER_SIZE), 
	}
}

func (r *trackerConnection) Connect() {
	conn, err := net.ListenMulticastUDP("udp", nil, r.addr)
	if err != nil {
		panic(err)
	}
	r.conn = conn
}

func (r *trackerConnection) Receive(latest *AtomicLatest[ssl_vision.TrackerWrapperPacket]) {
	for {
		sz, err := r.conn.Read(r.buff)
		if err != nil {
			Logger.Errorf("Unable to receive packet: %v", err)
			continue
		}

		pkt := new(ssl_vision.TrackerWrapperPacket) // heap allocate explicitly
		if err := proto.Unmarshal(r.buff[:sz], pkt); err != nil {
			Logger.Errorf("Unable to unmarshal packet: %v", err)
			continue
		}
		latest.Set(pkt) // publish latest atomically
	}
}

type TrackerSource struct {
	ssl    *trackerConnection
	latest *AtomicLatest[ssl_vision.TrackerWrapperPacket]
}

func NewTrackerSource(sslReceiverAddress string) *TrackerSource {
	client := &TrackerSource{
		ssl:    NewTrackerConnection(sslReceiverAddress),
		latest: NewAtomicLatest[ssl_vision.TrackerWrapperPacket](),
	}
	client.Connect()
	return client
}

func (receiver *TrackerSource) Connect() {
	receiver.ssl.Connect()
	go receiver.ssl.Receive(receiver.latest)
}

func (receiver *TrackerSource) GetMetaData() (string, string, float64, uint32, []ssl_vision.Capability) {
	pkt := receiver.latest.Get()
	if pkt == nil {
		return "", "", 0, 0, nil
	}
	frame := pkt.GetTrackedFrame()
	return pkt.GetSourceName(),
		pkt.GetUuid(),
		frame.GetTimestamp(),
		frame.GetFrameNumber(),
		frame.GetCapabilities()
}

func (receiver *TrackerSource) GetLatestPacket() (*ssl_vision.TrackerWrapperPacket, bool) {
	p := receiver.latest.Get()
	return p, p != nil
}

func (receiver *TrackerSource) GetTrackedBall() (info.BallState, uint32, bool) {
	pkt := receiver.latest.Get()
	if pkt == nil {
		return info.BallState{}, 0, false
	}
	frame := pkt.GetTrackedFrame()
	if frame == nil || len(frame.GetBalls()) == 0 {
		var fn uint32
		if frame != nil {
			fn = frame.GetFrameNumber()
		}
		return info.BallState{}, fn, false
	}

	ball := frame.GetBalls()[0] 
	pos := ball.GetPos()
	vel := ball.GetVel()

	state := info.NewBallState(
		info.Position{X: float64(pos.GetX() * 1000), Y: float64(pos.GetY() * 1000)},
		info.Position{X: float64(vel.GetX()), Y: float64(vel.GetY())},
		float64(ball.GetVisibility()),
		int64(frame.GetTimestamp()*1000),
		pkt.GetSourceName(),
	)
	return state, frame.GetFrameNumber(), true
}

func (receiver *TrackerSource) GetTrackedRobot(team info.Team, id info.ID) (info.RobotState, uint32, bool) {
	pkt := receiver.latest.Get()
	if pkt == nil {
		return info.RobotState{}, 0, false
	}
	frame := pkt.GetTrackedFrame()
	if frame == nil {
		return info.RobotState{}, 0, false
	}

	var wantTeam gc.Team
	switch team {
	case info.Yellow:
		wantTeam = gc.Team_YELLOW
	case info.Blue:
		wantTeam = gc.Team_BLUE
	default:
		wantTeam = gc.Team_UNKNOWN
	}

	for _, tr := range frame.GetRobots() {
		rid := tr.GetRobotId()
		if rid.GetTeam() != wantTeam || info.ID(rid.GetId()) != id {
			continue
		}

		pos := tr.GetPos()
		vel := tr.GetVel()
		state := info.NewRobotState(
			info.Position{
				X:     float64(pos.GetX() * 1000),
				Y:     float64(pos.GetY() * 1000),
				Z:     0,
				Angle: float64(tr.GetOrientation()),
			},
			info.Position{
				X: float64(vel.GetX()),
				Y: float64(vel.GetY()),
			},
			float64(tr.GetVisibility()),
			int64(frame.GetTimestamp()*1000),
			pkt.GetSourceName(),
		)
		return state, frame.GetFrameNumber(), true
	}

	return info.RobotState{}, frame.GetFrameNumber(), false
}

