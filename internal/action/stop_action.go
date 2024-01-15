package action

import (
	"github.com/LiU-SeeGoals/controller/internal/datatypes"
	"github.com/LiU-SeeGoals/proto-messages/robot_action"
)


type Stop struct {
	Id int
	isCommand bool
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