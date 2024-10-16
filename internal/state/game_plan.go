package state

import "fmt"

type RobotMove struct {
	Id       ID
	Position Position
}

type GamePlan struct {
	Valid        bool
	Team         Team
	Instructions []RobotMove
}

func NewGamePlan() GamePlan {
	gp := GamePlan{}
	return gp
}

func (gp GamePlan) ToDTO() string {
	dto := fmt.Sprintf("GamePlan{Valid: %t, Team: %d, Instructions: [", gp.Valid, gp.Team)
	for _, instruction := range gp.Instructions {
		dto += instruction.ToDTO() + ", "
	}

	return dto
}

func (rm RobotMove) ToDTO() string {
	dto := fmt.Sprintf("RobotMove{Id: %d, Position: %s}", rm.Id, rm.Position.ToDTO())
	return dto
}
