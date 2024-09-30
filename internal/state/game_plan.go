package state

type GamePlan struct {
	Valid bool
}

func NewGamePlan() *GamePlan {
	gp := &GamePlan{}
	return gp
}
