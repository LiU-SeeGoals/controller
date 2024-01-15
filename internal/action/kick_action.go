package action

import (
	"github.com/LiU-SeeGoals/controller/internal/datatypes"
	"github.com/LiU-SeeGoals/proto-messages/robot_action"
)

type Kick struct {
	Id int
	// 1 is slow, 10 is faster, limits unknown
	KickSpeed int
	isCommand bool
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