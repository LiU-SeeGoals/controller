package action

import (
	"github.com/LiU-SeeGoals/proto-messages/grsim"
)

type Placement interface {
	Translate() *grsim.GrSim_Replacement
	IsBall() bool
}

