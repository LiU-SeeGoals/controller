package action

import (
	"github.com/LiU-SeeGoals/controller/internal/proto/basestation"
	"gonum.org/v1/gonum/mat"
)

type RealAction int

const (
	RealStop RealAction = 0
	RealKick = 1
	RealMove = 2
	RealInit = 3
)

type Action interface {
	TranslateReal() *basestation.Command
	TranslateGrsim() int
}

type Stop struct {
	Id int
}

type Move struct {
	Id   int
	Pos  *mat.VecDense
	Goal *mat.VecDense
}

type Kick struct {
	Id int
	Speed int
}

type Init struct {
	Id int
}

func (m *Move) TranslateReal() *basestation.Command {
	command_move := &basestation.Command{
		CommandId: basestation.ActionType_MOVE_ACTION,
		RobotId: int32(m.Id),         
		Pos: &basestation.Vector3D{
			X: int32(m.Pos.AtVec(0)),
			Y: int32(m.Pos.AtVec(1)),
			W: float32(m.Pos.AtVec(2)),
		},
		Goal: &basestation.Vector3D{
			X: int32(m.Goal.AtVec(0)),
			Y: int32(m.Goal.AtVec(1)),
			W: float32(m.Goal.AtVec(2)),
		},
	}

	return command_move
}

func (m *Move) TranslateGrsim() int {
	return 0
}

func (s *Stop) TranslateReal() *basestation.Command {
	command_move := &basestation.Command{
		CommandId: basestation.ActionType_STOP_ACTION,
		RobotId: int32(s.Id),
	}

	return command_move
}

func (s *Stop) TranslateGrsim() int {
	return 0
}

func (k *Kick) TranslateReal() *basestation.Command {
	command_move := &basestation.Command{
		CommandId: basestation.ActionType_KICK_ACTION,
		RobotId: int32(k.Id),
		Speed: int32(k.Speed),
	}

	return command_move
}

func (s *Kick) TranslateGrsim() int {
	return 0
}

func (i *Init) TranslateReal() *basestation.Command {

	command_move := &basestation.Command{
		CommandId: basestation.ActionType_INIT_ACTION,
		RobotId: int32(i.Id),
	}

	return command_move
}