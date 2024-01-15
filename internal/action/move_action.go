package action

import (
	"math"

	"github.com/LiU-SeeGoals/controller/internal/gamestate/robot"
	"github.com/LiU-SeeGoals/proto-messages/grsim"
	"github.com/LiU-SeeGoals/proto-messages/robot_action"
	"gonum.org/v1/gonum/mat"
	"google.golang.org/protobuf/proto"
)

// Move represents an action for a robot on the field. It contains all the necessary
// information for moving a robot, including its ID, team, current position, destination,
// and whether it should dribble. This struct implements the Action interface.
type Move struct {
    Id        int         	// Id of the robot.
    Team      robot.Team  	// Team of the robot (Blue or Yellow).
    Pos       *mat.VecDense // Position {x (mm), y (mm), angle (rad)}.
    Dest      *mat.VecDense // Destination {x (mm), y (mm), angle (rad)}.
    Dribble   bool        	// Flag to activate the dribbler.
}

// NewMove creates a new Move action for a robot.
// id is the unique identifier for the robot. team specifies the robot's team.
// pos and dest are VecDense pointers representing the robot's current and destination positions, 
// each containing X, Y coordinates (in mm) and an orientation angle (in radians).
// dribble indicates whether the robot should activate its dribbler.
func NewMove(id int, team robot.Team, pos, dest *mat.VecDense, dribble bool) *Move {
	return &Move{
		Id: id,
		Team: team,
		Pos: pos,
		Dest: dest,
		Dribble: dribble,
	}
}

// TranslateGrsim translates the Move action into a command suitable for the grSim simulator.
// It calculates the required tangential and angular velocities and creates a GrSim_Robot_Command 
// struct which can be sent to the grSim simulator.
func (m *Move) TranslateGrsim() *grsim.GrSim_Robot_Command {

	Veltangent, VelAngular := m.calculateVelocities()

	return &grsim.GrSim_Robot_Command{
		Id: proto.Uint32(uint32(m.Id)),
		Kickspeedx: proto.Float32(0),
		Kickspeedz: proto.Float32(0),
		Veltangent: proto.Float32(Veltangent),
		Velnormal: proto.Float32(0),
		Velangular: proto.Float32(VelAngular),
		Spinner: &m.Dribble,
		Wheelsspeed: proto.Bool(false),
	}
}

// TranslateReal translates the Move action into a real-world command.
// It converts the current and destination positions into a robot_action.Command format.
func (m *Move) TranslateReal() *robot_action.Command {
	command_move := &robot_action.Command{
		CommandId: robot_action.ActionType_MOVE_ACTION,
		RobotId:   int32(m.Id),
		Pos: VecDenseToVector3D(m.Pos),
		Dest: VecDenseToVector3D(m.Dest),
	}

	return command_move
}

// IsTeamYellow checks if the robot's team is Yellow.
func (m *Move) IsTeamYellow() bool {
	return m.Team == robot.Yellow
}

// calculateVelocities calculates the tangential and angular velocities needed to move the robot
// from its current position to its destination. This method is used internally by TranslateGrsim.
// It returns the calculated tangential and angular velocities.
func (m *Move) calculateVelocities() (float32, float32) {

	var VelAngular float32 = 0
	var VelTangent float32 = 0

	diff := mat.NewVecDense(3, nil)
	diff.SubVec(m.Dest, m.Pos)

	angle := math.Atan2(diff.AtVec(1), diff.AtVec(0))
	diffPosAngle := angle - m.Pos.AtVec(2)
	diffDestAngle := m.Pos.AtVec(2) - m.Dest.AtVec(2)

	if math.Abs(diff.AtVec(0)) > 50 || math.Abs(diff.AtVec(1)) > 50 {

		if diffPosAngle > 0.2 {
			VelAngular = 4
		} else if diffPosAngle < -0.2 {
			VelAngular = -4
		} else {
			VelTangent = 1
		}
	} else if diffDestAngle > 0.2 {
		VelAngular = -4
	} else if diffDestAngle < -0.2 {
		VelAngular = 4
	}

	return VelTangent, VelAngular
}