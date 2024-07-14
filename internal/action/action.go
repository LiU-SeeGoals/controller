package action

import (
	"math"

	"github.com/LiU-SeeGoals/proto_go/robot_action"
	"github.com/LiU-SeeGoals/proto_go/simulation"
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
	// Translates an action to parameters defined for sim
	TranslateSim() *simulation.RobotCommand
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

type MoveTo struct {
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
type Move struct {
	Id        int
	Direction *mat.VecDense // 2D vector, first value is x, second is y
}

type Init struct {
	Id int
}

type Ping struct {
    Id int
}

//------------------------------------------------------------------//
// TranslateSim translates the action to simulation proto message	//
// (there are a lot of wrapper proto messages :(                    //
//------------------------------------------------------------------//

func (s *Stop) TranslateSim() *simulation.RobotCommand {
	id := uint32(s.Id)
	angular := float32(0)
	forward := float32(0)
	left := float32(0)

	localVel := &simulation.MoveLocalVelocity{
		Forward: &forward,
		Left:    &left,
		Angular: &angular,
	}

	moveCommand := &simulation.RobotMoveCommand{
		Command: &simulation.RobotMoveCommand_LocalVelocity{
			LocalVelocity: localVel,
		},
	}

	return &simulation.RobotCommand{
		Id:          &id,
		MoveCommand: moveCommand,
	}

}

func (mv *MoveTo) TranslateSim() *simulation.RobotCommand {

	id := uint32(mv.Id)
	diff := mat.NewVecDense(3, nil)
	diff.SubVec(mv.Dest, mv.Pos)

	dribblerSpeed := float32(0)
	if mv.Dribble {
		dribblerSpeed = 100 // in rpm, adjust as needed
	}

	goalAngle := math.Atan2(diff.AtVec(1), diff.AtVec(0))
	currAngle := mv.Pos.AtVec(2)
	inPosition := false

	if math.Abs(diff.AtVec(0)) < 50 && math.Abs(diff.AtVec(1)) < 50 {
		inPosition = true
		goalAngle = mv.Dest.AtVec(2)

	}

	// Normalize an angle to be within -π to π
	normalizeAngle := func(angle float64) float64 {
		angle = math.Mod(angle+math.Pi, 2*math.Pi)
		if angle < 0 {
			angle += 2 * math.Pi
		}
		return angle - math.Pi
	}

	// Calculate difference between current angle and goal angle
	angleDiff := normalizeAngle(goalAngle - currAngle)

	// Set angular velocity
	angular := float32(0)
	if angleDiff > 0.2 {
		angular = 2 // Adjust this value as necessary
	} else if angleDiff < -0.2 {
		angular = -2 // Adjust this value as necessary
	}

	// Set forward speed
	forward := float32(0)
	if math.Abs(angleDiff) < 0.3 && !inPosition {
		forward = 1 // Move forward when facing the goal
	}

	left := float32(0)

	// Create the local velocity command
	localVel := &simulation.MoveLocalVelocity{
		Forward: &forward,
		Left:    &left,
		Angular: &angular,
	}

	// Create the move command and assign the local velocity to the oneof field
	moveCommand := &simulation.RobotMoveCommand{
		Command: &simulation.RobotMoveCommand_LocalVelocity{
			LocalVelocity: localVel,
		},
	}

	// Create the robot command with the move command
	return &simulation.RobotCommand{
		Id:            &id,
		MoveCommand:   moveCommand,
		DribblerSpeed: &dribblerSpeed,
	}
}

func (d *Dribble) TranslateSim() *simulation.RobotCommand {
	id := uint32(d.Id)
	dribblerSpeed := float32(0)
	if d.Dribble {
		dribblerSpeed = 100 // in rpm, adjust as needed
	}

	return &simulation.RobotCommand{
		Id:            &id,
		DribblerSpeed: &dribblerSpeed,
	}
}

func (k *Kick) TranslateSim() *simulation.RobotCommand {
	id := uint32(k.Id)
	kickSpeed := float32(k.KickSpeed) // in m/s

	return &simulation.RobotCommand{
		Id:        &id,
		KickSpeed: &kickSpeed,
	}
}

func (r *Rotate) TranslateSim() *simulation.RobotCommand {
	id := uint32(r.Id)
	angular := float32(r.AngularVel) // No angular velocity currently, adjust as needed
	forward := float32(0)
	left := float32(0)

	localVel := &simulation.MoveLocalVelocity{
		Forward: &forward,
		Left:    &left,
		Angular: &angular,
	}

	moveCommand := &simulation.RobotMoveCommand{
		Command: &simulation.RobotMoveCommand_LocalVelocity{
			LocalVelocity: localVel,
		},
	}

	return &simulation.RobotCommand{
		Id:          &id,
		MoveCommand: moveCommand,
	}
}

func (s *Move) TranslateSim() *simulation.RobotCommand {

	id := uint32(s.Id)
	angular := float32(0) // No angular velocity currently, adjust as needed
	forward := float32(s.Direction.AtVec(0))
	left := float32(s.Direction.AtVec(1))

	// Create the local velocity command
	localVel := &simulation.MoveLocalVelocity{
		Forward: &forward,
		Left:    &left,
		Angular: &angular,
	}

	// Create the move command and assign the local velocity to the oneof field
	moveCommand := &simulation.RobotMoveCommand{
		Command: &simulation.RobotMoveCommand_LocalVelocity{
			LocalVelocity: localVel,
		},
	}

	// Create the robot command with the move command
	return &simulation.RobotCommand{
		Id:          &id,
		MoveCommand: moveCommand,
	}
}

// Do nothing, only implemented to satisfy interface
func (i *Init) TranslateSim() *simulation.RobotCommand {
	id := uint32(i.Id)
	return &simulation.RobotCommand{
		Id: &id,
	}
}

// Do nothing, only implemented to satisfy interface
func (i *Ping) TranslateSim() *simulation.RobotCommand {
	id := uint32(i.Id)
	return &simulation.RobotCommand{
		Id: &id,
	}
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

func (mt *MoveTo) TranslateReal() *robot_action.Command {
	command_move := &robot_action.Command{
		CommandId: robot_action.ActionType_MOVE_TO_ACTION,
		RobotId:   int32(mt.Id),
		Pos: &robot_action.Vector3D{
			X: int32(mt.Pos.AtVec(0)),
			Y: int32(mt.Pos.AtVec(1)),
			W: float32(mt.Pos.AtVec(2)),
		},
		Dest: &robot_action.Vector3D{
			X: int32(mt.Dest.AtVec(0)),
			Y: int32(mt.Dest.AtVec(1)),
			W: float32(mt.Dest.AtVec(2)),
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

func (s *Move) TranslateReal() *robot_action.Command {
	command := &robot_action.Command{
		CommandId: robot_action.ActionType_MOVE_ACTION,
		RobotId:   int32(s.Id),
		Direction: &robot_action.Vector2D{
			X: int32(s.Direction.AtVec(0)),
			Y: int32(s.Direction.AtVec(1)),
		},
	}
	return command
}

func (s *Ping) TranslateReal() *robot_action.Command {
	command := &robot_action.Command{
		CommandId: robot_action.ActionType_MOVE_ACTION,
		RobotId:   int32(s.Id),
	}
	return command
}

//----------------------------------------------------------------------------------------------
// ToDTO
//----------------------------------------------------------------------------------------------

func (m *MoveTo) ToDTO() ActionDTO {
	return ActionDTO{
		Action:  robot_action.ActionType_MOVE_TO_ACTION,
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

func (s *Move) ToDTO() ActionDTO {
	return ActionDTO{
		Action: robot_action.ActionType_MOVE_ACTION,
		Id:     s.Id,
		DestX:  int32(s.Direction.AtVec(0)),
		DestY:  int32(s.Direction.AtVec(1)),
		DestW:  0,
	}
}

func (s *Ping) ToDTO() ActionDTO {
	return ActionDTO{
		Action: robot_action.ActionType_PING,
		Id:     s.Id,
	}
}
