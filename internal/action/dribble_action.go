package action

import "github.com/LiU-SeeGoals/controller/internal/datatypes"

type Dribble struct {
	Id int
	// set Dribbling, useless right now
	Dribble bool
	isCommand bool
}

func (d *Dribble) TranslateGrsim(params *datatypes.Parameters) {
	params.RobotId = uint32(d.Id)
	params.Spinner = d.Dribble
}