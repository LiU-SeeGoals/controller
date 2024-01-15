package action

import (
	"github.com/LiU-SeeGoals/controller/internal/datatypes"
	"github.com/LiU-SeeGoals/proto-messages/robot_action"
)

type Init struct {
	Id int
	isCommand bool
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