package action

import (
	"github.com/LiU-SeeGoals/controller/internal/datatypes"
	"github.com/LiU-SeeGoals/proto-messages/robot_action"
	"gonum.org/v1/gonum/mat"
)

// Forward is x=0, y=1, Backward is x=0, y=-1, Left is x=-1, y=0, Right is x=1, y=0
// the size of the vector sets the speed of the robot
type SetNavigationDirection struct {
	Id        int
	Direction *mat.VecDense // 2D vector, first value is x, second is y
	isCommand bool
}

func (s *SetNavigationDirection) TranslateGrsim(params *datatypes.Parameters) {
	params.RobotId = uint32(s.Id)
	params.VelNormal = float32(s.Direction.AtVec(0))
	params.VelTangent = float32(s.Direction.AtVec(1))
}

func (s *SetNavigationDirection) TranslateReal() *robot_action.Command {
	command := &robot_action.Command{
		CommandId: robot_action.ActionType_SET_NAVIGATION_DIRECTION_ACTION,
		RobotId:   int32(s.Id),
		Direction: &robot_action.Vector2D{
			X: int32(s.Direction.AtVec(0)),
			Y: int32(s.Direction.AtVec(1)),
		},
	}

	return command
}