package action

import (
	"github.com/LiU-SeeGoals/controller/internal/gamestate/robot"
	"github.com/LiU-SeeGoals/proto-messages/grsim"
	"gonum.org/v1/gonum/mat"
	"google.golang.org/protobuf/proto"
)

type RobotPlacement struct {
	Id int
	Team robot.Team
	Pos *mat.VecDense
	TurnOn bool
}

func newRobotPlacement(id int, team robot.Team, pos *mat.VecDense, turnOn bool) *RobotPlacement {
	return &RobotPlacement{
		Id: id,
		Team: team,
		Pos: pos,
		TurnOn: turnOn,
	}
}

func (r *RobotPlacement) Translate() *grsim.GrSim_Replacement {
	robotReplacement := &grsim.GrSim_RobotReplacement{
		Id: proto.Uint32(uint32(r.Id)),
		X: proto.Float64(r.Pos.AtVec(0)),
		Y: proto.Float64(r.Pos.AtVec(1)),
		Dir: proto.Float64(r.Pos.AtVec(2)),
		Yellowteam: proto.Bool(r.Team == robot.Yellow),
		Turnon: proto.Bool(r.TurnOn),
	}

	return &grsim.GrSim_Replacement{
		Robots: []*grsim.GrSim_RobotReplacement{robotReplacement},
	}
}

func (p *RobotPlacement) isBall() bool {
	return false
}
