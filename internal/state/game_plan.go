package state

type RobotMove struct {
	Id       ID
	Position Position
}

type GamePlan struct {
	Valid        bool
	Team         Team
	Instructions []RobotMove
}

func NewGamePlan() *GamePlan {
	gp := &GamePlan{}
	return gp
}
