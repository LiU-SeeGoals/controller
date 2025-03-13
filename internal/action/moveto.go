package action

import (
	"math"

	"github.com/LiU-SeeGoals/proto_go/robot_action"
	"github.com/LiU-SeeGoals/proto_go/simulation"

	"github.com/LiU-SeeGoals/controller/internal/info"
)

type MoveTo struct {
	// The id of the robot.
	Id int
	// Current position of Robot, vector contains (x,y,w)
	Pos info.Position
	// Goal destination of Robot, vector contains (x,y,w)
	Dest info.Position
	// Decides if the robot should dribble while moving
	Dribble bool
	// We need to know ID AND team to know how to update the pos
	Team info.Team
}

func (mv *MoveTo) TranslateSim() *simulation.RobotCommand {
	id := uint32(mv.Id)

	// Angular velocity counter-clockwise [rad/s]
	dx := float64(mv.Pos.X - mv.Dest.X)
	dy := float64(mv.Pos.Y - mv.Dest.Y)
	angleDiff := mv.Dest.Angle - mv.Pos.Angle

	if angleDiff > math.Pi {
		angleDiff -= 2 * math.Pi
	}
	if angleDiff < -math.Pi {
		angleDiff += 2 * math.Pi
	}

	distance := math.Sqrt(dx*dx + dy*dy)
	maxSpeed := float64(0.5)
	DeAccDistance := float64(300) // The distance from target robot start to deaccelerate (measured in mm)
	speed := float32(math.Min(maxSpeed, (maxSpeed/DeAccDistance)*distance))

	maxAngleSpeed := float64(2)
	deAccAngleDistance := float64(0.5) // The distance from target robot start to deaccelerate (measured in rad)
	angle := float32(math.Min(maxAngleSpeed, (maxAngleSpeed/deAccAngleDistance)*float64(angleDiff)))

	// Compute the target direction in global space
	targetDirection := math.Atan2(-dy, -dx)
	targetDirection = math.Mod(targetDirection+math.Pi, 2*math.Pi) - math.Pi

	moveAngle := targetDirection - mv.Pos.Angle

	// Decompose movement into forward and leftward velocities
	forward := speed * float32(math.Cos(moveAngle)) // Forward velocity
	left := speed * float32(math.Sin(moveAngle))    // Leftward velocity

	dribblerSpeed := float32(0)
	if mv.Dribble {
		dribblerSpeed = 100 // in rpm, adjust as needed
	}

	localVel := &simulation.MoveLocalVelocity{
		Forward: &forward,
		Left:    &left,
		Angular: &angle,
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

func (mt *MoveTo) TranslateReal() *robot_action.Command {
	command_move := &robot_action.Command{
		CommandId: robot_action.ActionType_MOVE_TO_ACTION,
		RobotId:   int32(mt.Id),
		Pos: &robot_action.Vector3D{
			X: int32(mt.Pos.X + 10000),
			Y: int32(mt.Pos.Y + 10000),
			W: float32(mt.Pos.Angle * 1000),
		},
		Dest: &robot_action.Vector3D{
			X: int32(mt.Dest.X + 10000),
			Y: int32(mt.Dest.Y + 10000),
			W: float32(mt.Dest.Angle * 1000),
		},
	}
	return command_move
}

func (m *MoveTo) ToDTO() ActionDTO {
	return ActionDTO{
		Action:  robot_action.ActionType_MOVE_TO_ACTION,
		Id:      m.Id,
		PosX:    int32(m.Pos.X),
		PosY:    int32(m.Pos.Y),
		PosW:    float32(m.Pos.Angle),
		DestX:   int32(m.Dest.X),
		DestY:   int32(m.Dest.Y),
		DestW:   float32(m.Dest.Angle),
		Dribble: m.Dribble,
	}
}
