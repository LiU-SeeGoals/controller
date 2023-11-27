package action

import (
	"fmt"
	"math"

	"github.com/LiU-SeeGoals/controller/internal/datatypes"
	"gonum.org/v1/gonum/mat"
)

type Action interface {
	//TranslateReal() int
	TranslateGrsim(params *datatypes.Parameters)
}

type Stop struct {
}

type Move struct {
	Pos     *mat.VecDense
	Dest    *mat.VecDense
	Dribble bool
}

type Dribble struct {
	Dribble bool
}

type Kick struct {
	Kickspeed int
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
