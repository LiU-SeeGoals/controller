package action

import (
	"math"

	"github.com/LiU-SeeGoals/controller/internal/datatypes"
	"gonum.org/v1/gonum/mat"
)

type Action interface {
	//TranslateReal() int
	TranslateGrsim(params *datatypes.Parameters)
}

type Stop struct{}

type GoTo struct {
	Pos  *mat.VecDense
	Dest *mat.VecDense
}

func (gt *GoTo) TranslateGrsim(params *datatypes.Parameters) {
	diff := mat.NewVecDense(3, nil)
	diff.SubVec(gt.Dest, gt.Pos)

	if diff.AtVec(0) < 1 && diff.AtVec(1) < 1 {
	} else {
		angle := math.Atan2(diff.AtVec(1), diff.AtVec(0))
		if angle < 2 {
			params.VelNormal = 50
		} else {

			params.VelAngular = 5
		}
	}

}

// type Pass struct {
// }

// func (s Stop) TranslateGrsim() int {

// }
