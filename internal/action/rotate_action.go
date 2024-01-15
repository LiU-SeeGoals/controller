package action

import (
	"github.com/LiU-SeeGoals/controller/internal/datatypes"
	"github.com/LiU-SeeGoals/proto-messages/robot_action"
)

// Negative value rotates robot clockwise
type Rotate struct {
	Id         int
	AngularVel int
	isCommand bool
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