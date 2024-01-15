package action

import (
	"github.com/LiU-SeeGoals/proto-messages/grsim"
	"github.com/LiU-SeeGoals/proto-messages/robot_action"
	"gonum.org/v1/gonum/mat"
)

type Action interface {
	TranslateReal() *robot_action.Command
	// Translates an action to parameters defined for grsim
	TranslateGrsim() *grsim.GrSim_Robot_Command
	IsTeamYellow() bool
}

// VecDenseToVector3D converts a gonum/mat VecDense to a robot_action Vector3D.
// It is used to translate vector formats between gonum/mat and the protobuf message.
func VecDenseToVector3D(vec *mat.VecDense) *robot_action.Vector3D {
	return &robot_action.Vector3D{
		X: int32(vec.AtVec(0)),
		Y: int32(vec.AtVec(1)),
		W: float32(vec.AtVec(2)),
	}
}
