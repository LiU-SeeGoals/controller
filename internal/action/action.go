package action

import (
	"math"

	"github.com/LiU-SeeGoals/controller/internal/datatypes"
	// "github.com/LiU-SeeGoals/controller/internal/proto/basestation"
	"github.com/LiU-SeeGoals/proto_go/robot_action"
	"gonum.org/v1/gonum/mat"
)

// import (
// 	"math"

// 	"github.com/LiU-SeeGoals/controller/internal/datatypes"
// 	"github.com/LiU-SeeGoals/controller/internal/proto/basestation"
// 	"gonum.org/v1/gonum/mat"
// )

type Action interface {
	TranslateReal() *robot_action.Command
	// Translates an action to parameters defined for grsim
	TranslateGrsim(params *datatypes.Parameters)
	ToDTO() ActionDTO
}

type ActionDTO struct {
	// The id of the robot.
	Id     int                     `json:"Id"`
	Action robot_action.ActionType `json:"Action"`
	// Current position of Robot, vector contains (x,y,w)
	PosX int32   `json:"PosX"`
	PosY int32   `json:"PosY"`
	PosW float32 `json:"PosW"`
	// Goal destination of Robot, vector contains (x,y,w)
	DestX int32   `json:"DestX"`
	DestY int32   `json:"DestY"`
	DestW float32 `json:"DestW"`
	// Decides if the robot should dribble while moving
	Dribble bool `json:"Dribble"`
}

//----------------------------------------------------------------------------------------------
// Actions structs
//----------------------------------------------------------------------------------------------

type Stop struct {
	Id int
}

type Move struct {
	// The id of the robot.
	Id int
	// Current position of Robot, vector contains (x,y,w)
	Pos *mat.VecDense
	// Goal destination of Robot, vector contains (x,y,w)
	Dest *mat.VecDense
	// Decides if the robot should dribble while moving
	Dribble bool
}

type Dribble struct {
	Id int
	// set Dribbling, useless right now
	Dribble bool
}

type Kick struct {
	Id int
	// 1 is slow, 10 is faster, limits unknown
	KickSpeed int
}

// Negative value rotates robot clockwise
type Rotate struct {
	Id         int
	AngularVel int
}

// Forward is x=0, y=1, Backward is x=0, y=-1, Left is x=-1, y=0, Right is x=1, y=0
// the size of the vector sets the speed of the robot
type SetNavigationDirection struct {
	Id        int
	Direction *mat.VecDense // 2D vector, first value is x, second is y
}

type Init struct {
	Id int
}

//----------------------------------------------------------------------------------------------
// TranslateGrsim
//----------------------------------------------------------------------------------------------

func (s *SetNavigationDirection) TranslateGrsim(params *datatypes.Parameters) {
	params.RobotId = uint32(s.Id)
	params.VelNormal = float32(s.Direction.AtVec(0))
	params.VelTangent = float32(s.Direction.AtVec(1))
}

func (r *Rotate) TranslateGrsim(params *datatypes.Parameters) {
	params.RobotId = uint32(r.Id)
	params.VelAngular = float32(r.AngularVel)
}

func (k *Kick) TranslateGrsim(params *datatypes.Parameters) {
	params.RobotId = uint32(k.Id)
	params.KickSpeedX = float32(k.KickSpeed)
}

func (s *Stop) TranslateGrsim(params *datatypes.Parameters) {
	params.RobotId = uint32(s.Id)
	params.VelNormal = float32(0)
	params.VelTangent = float32(0)
}

func (mv *Move) TranslateGrsim(params *datatypes.Parameters) {
	params.RobotId = uint32(mv.Id)
	diff := mat.NewVecDense(3, nil)
	diff.SubVec(mv.Dest, mv.Pos)
	params.Spinner = mv.Dribble

	angle := math.Atan2(diff.AtVec(1), diff.AtVec(0))
	diffPosAngle := angle - mv.Pos.AtVec(2)
	diffDestAngle := mv.Pos.AtVec(2) - mv.Dest.AtVec(2)

	if math.Abs(diff.AtVec(0)) > 50 || math.Abs(diff.AtVec(1)) > 50 {

		if diffPosAngle > 0.2 {
			params.VelAngular = 4
		} else if diffPosAngle < -0.2 {
			params.VelAngular = -4
		} else {
			params.VelTangent = 5
		}
	} else if diffDestAngle > 0.2 {
		params.VelAngular = -4
	} else if diffDestAngle < -0.2 {
		params.VelAngular = 4
	}
}

func (d *Dribble) TranslateGrsim(params *datatypes.Parameters) {
	params.RobotId = uint32(d.Id)
	params.Spinner = d.Dribble
}

func (i *Init) TranslateGrsim(params *datatypes.Parameters) {
	params.RobotId = uint32(i.Id)
}

//----------------------------------------------------------------------------------------------
// TranslateReal
//----------------------------------------------------------------------------------------------

func (r *Rotate) TranslateReal() *robot_action.Command {
	command_move := &robot_action.Command{
		CommandId:  robot_action.ActionType_ROTATE_ACTION,
		RobotId:    int32(r.Id),
		AngularVel: int32(r.AngularVel),
	}
	return command_move
}

func (k *Kick) TranslateReal() *robot_action.Command {
	command_move := &robot_action.Command{
		CommandId: robot_action.ActionType_KICK_ACTION,
		RobotId:   int32(k.Id),
		KickSpeed: int32(k.KickSpeed),
	}
	return command_move
}

func (s *Stop) TranslateReal() *robot_action.Command {
	command_move := &robot_action.Command{
		CommandId: robot_action.ActionType_STOP_ACTION,
		RobotId:   int32(s.Id),
	}
	return command_move
}

func (m *Move) TranslateReal() *robot_action.Command {
	command_move := &robot_action.Command{
		CommandId: robot_action.ActionType_MOVE_ACTION,
		RobotId:   int32(m.Id),
		Pos: &robot_action.Vector3D{
			X: int32(m.Pos.AtVec(0)),
			Y: int32(m.Pos.AtVec(1)),
			W: float32(m.Pos.AtVec(2)),
		},
		Dest: &robot_action.Vector3D{
			X: int32(m.Dest.AtVec(0)),
			Y: int32(m.Dest.AtVec(1)),
			W: float32(m.Dest.AtVec(2)),
		},
	}
	return command_move
}

func (i *Init) TranslateReal() *robot_action.Command {

	command_move := &robot_action.Command{
		CommandId: robot_action.ActionType_INIT_ACTION,
		RobotId:   int32(i.Id),
	}
	return command_move
}

func (s *SetNavigationDirection) TranslateReal() *robot_action.Command {
	command := &robot_action.Command{
		CommandId: robot_action.ActionType_MOVE_TO_ACTION,
		RobotId:   int32(s.Id),
		Direction: &robot_action.Vector2D{
			X: int32(s.Direction.AtVec(0)),
			Y: int32(s.Direction.AtVec(1)),
		},
	}
	return command
}

//----------------------------------------------------------------------------------------------
// ToDTO
//----------------------------------------------------------------------------------------------

func (m *Move) ToDTO() ActionDTO {
	return ActionDTO{
		Action:  robot_action.ActionType_MOVE_ACTION,
		Id:      m.Id,
		PosX:    int32(m.Pos.AtVec(0)),
		PosY:    int32(m.Pos.AtVec(1)),
		PosW:    float32(m.Pos.AtVec(2)),
		DestX:   int32(m.Dest.AtVec(0)),
		DestY:   int32(m.Dest.AtVec(1)),
		DestW:   float32(m.Dest.AtVec(2)),
		Dribble: m.Dribble,
	}
}

func (i *Init) ToDTO() ActionDTO {
	return ActionDTO{
		Action: robot_action.ActionType_INIT_ACTION,
		Id:     i.Id,
	}
}

func (r *Rotate) ToDTO() ActionDTO {
	return ActionDTO{
		Action: robot_action.ActionType_ROTATE_ACTION,
		Id:     r.Id,
		PosW:   float32(r.AngularVel),
	}
}

func (k *Kick) ToDTO() ActionDTO {
	return ActionDTO{
		Action: robot_action.ActionType_KICK_ACTION,
		Id:     k.Id,
		PosW:   float32(k.KickSpeed),
	}
}

func (s *Stop) ToDTO() ActionDTO {
	return ActionDTO{
		Action: robot_action.ActionType_STOP_ACTION,
		Id:     s.Id,
	}
}

func (s *SetNavigationDirection) ToDTO() ActionDTO {
	return ActionDTO{
		Action: robot_action.ActionType_MOVE_TO_ACTION,
		Id:     s.Id,
		DestX:  int32(s.Direction.AtVec(0)),
		DestY:  int32(s.Direction.AtVec(1)),
		DestW:  0,
	}
}
