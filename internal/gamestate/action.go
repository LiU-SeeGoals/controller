package action

import (
	"gonum.org/v1/gonum/mat"
)

type Action interface {
	TranslateReal() int
	TranslateGrsim() int
}

type Stop struct{}

type GoTo struct {
	pos *mat.VecDense
}

type Kick struct {
}

type Init struct {
}

type Pass struct {
}
