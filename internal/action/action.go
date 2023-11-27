package action

import (
	"fmt"
	"math"

	"github.com/LiU-SeeGoals/controller/internal/datatypes"
	"gonum.org/v1/gonum/mat"
)

type Action interface {
	// TranslateReal() int
	// Translates an action to parameters defined for grsim
	TranslateGrsim(params *datatypes.Parameters)
}

type Stop struct {
}

type Move struct {
	// Current position of Robot, vector contains (x,y,w)
	Pos *mat.VecDense
	// Goal destination of Robot, vector contains (x,y,w)
	Dest *mat.VecDense
	// Decides if the robot should dribble while moving
	Dribble bool
}

type Dribble struct {
	// set Dribbling, useless right now
	Dribble bool
}

type Kick struct {
	// 1 is slow, 10 is faster, limits unknown
	Kickspeed int
}

// Negative value rotates robot clockwise
type Rotate struct {
	AngularVel int
}

// index 0: positive left
// index 1: positive forward
type SetNavigationDirection struct {
	Direction *mat.VecDense
}

func (s *SetNavigationDirection) TranslateGrsim(params *datatypes.Parameters) {
	params.VelNormal = float32(s.Direction.AtVec(0))
	params.VelTangent = float32(s.Direction.AtVec(1))
}

func (r *Rotate) TranslateGrsim(params *datatypes.Parameters) {
	params.VelAngular = float32(r.AngularVel)
}

func (d *Dribble) TranslateGrsim(params *datatypes.Parameters) {

	params.Spinner = d.Dribble
}

func (k *Kick) TranslateGrsim(params *datatypes.Parameters) {

	params.KickSpeedX = float32(k.Kickspeed)
}

func (s *Stop) TranslateGrsim(params *datatypes.Parameters) {

}

func (mv *Move) TranslateGrsim(params *datatypes.Parameters) {
	diff := mat.NewVecDense(3, nil)
	diff.SubVec(mv.Dest, mv.Pos)
	params.Spinner = mv.Dribble

	angle := math.Atan2(diff.AtVec(1), diff.AtVec(0))
	diffPosAngle := angle - mv.Pos.AtVec(2)
	diffDestAngle := mv.Pos.AtVec(2) - mv.Dest.AtVec(2)

	fmt.Println(diffPosAngle)
	if math.Abs(diff.AtVec(0)) > 100 || math.Abs(diff.AtVec(1)) > 100 {

		if diffPosAngle > 0.2 {
			params.VelAngular = 4
		} else if diffPosAngle < -0.2 {
			params.VelAngular = -4
		} else {
			params.VelTangent = 1
		}
	} else if diffDestAngle > 0.2 {
		params.VelAngular = -4
	} else if diffDestAngle < -0.2 {
		params.VelAngular = 4
	}

}
