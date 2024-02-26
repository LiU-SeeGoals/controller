package action

import (
	"math"

	"github.com/LiU-SeeGoals/controller/internal/datatypes"
	"github.com/LiU-SeeGoals/proto_go/robot_action"
	"gonum.org/v1/gonum/mat"
)

type Action interface {
	TranslateReal() *robot_action.Command
	// Translates an action to parameters defined for grsim
	TranslateGrsim(params *datatypes.Parameters)
}

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

func (s *Move) TranslateGrsim(params *datatypes.Parameters) {
	params.RobotId = uint32(s.Id)
	params.VelNormal = float32(s.Direction.AtVec(0))
	params.VelTangent = float32(s.Direction.AtVec(1))
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

func (r *Rotate) TranslateGrsim(params *datatypes.Parameters) {
	params.RobotId = uint32(r.Id)
	params.VelAngular = float32(r.AngularVel)
}

func (r *Rotate) TranslateReal() *robot_action.Command {
	command_move := &robot_action.Command{
		CommandId:  robot_action.ActionType_ROTATE_ACTION,
		RobotId:    int32(r.Id),
		AngularVel: int32(r.AngularVel),
	}

	return command_move
}

func (d *Dribble) TranslateGrsim(params *datatypes.Parameters) {
	params.RobotId = uint32(d.Id)
	params.Spinner = d.Dribble
}

func (k *Kick) TranslateGrsim(params *datatypes.Parameters) {
	params.RobotId = uint32(k.Id)
	params.KickSpeedX = float32(k.KickSpeed)
}

func (k *Kick) TranslateReal() *robot_action.Command {
	command_move := &robot_action.Command{
		CommandId: robot_action.ActionType_KICK_ACTION,
		RobotId:   int32(k.Id),
		KickSpeed: int32(k.KickSpeed),
	}

	return command_move
}

func (s *Stop) TranslateGrsim(params *datatypes.Parameters) {

	params.RobotId = uint32(s.Id)
	params.VelNormal = float32(0)
	params.VelTangent = float32(0)

}

func (s *Stop) TranslateReal() *robot_action.Command {
	command_move := &robot_action.Command{
		CommandId: robot_action.ActionType_STOP_ACTION,
		RobotId:   int32(s.Id),
	}

	return command_move
}

func (mv *MoveTo) TranslateGrsim(params *datatypes.Parameters) {

	params.RobotId = uint32(mv.Id)
	diff := mat.NewVecDense(3, nil)
	diff.SubVec(mv.Dest, mv.Pos)
	params.Spinner = mv.Dribble

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
	if angleDiff > 0.2 {
		params.VelAngular = 2 // Adjust this value as necessary
	} else if angleDiff < -0.2 {
		params.VelAngular = -2 // Adjust this value as necessary
	}

	// Set forward speed
	if math.Abs(angleDiff) < 0.3 && !inPosition {
		params.VelTangent = 1 // Move forward when facing the goal
	}
}

func (m *MoveTo) TranslateReal() *robot_action.Command {
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

func (i *Init) TranslateGrsim(params *datatypes.Parameters) {
	params.RobotId = uint32(i.Id)
}

func (i *Init) TranslateReal() *robot_action.Command {

	command_move := &robot_action.Command{
		CommandId: robot_action.ActionType_INIT_ACTION,
		RobotId:   int32(i.Id),
	}

	return command_move
}
