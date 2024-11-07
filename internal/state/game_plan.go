package state

import "fmt"

type InstructionType int

const (
	MoveToPosition         InstructionType = 0 // Move to a specific position (static)
	MoveToBall             InstructionType = 1 // Move to the ball (dynamic)
	MoveToBallFacePosition InstructionType = 2 // Move to the ball and face a specific position (static)
	MoveToBallFacePlayer   InstructionType = 3 // Move to the ball and face a specific player (dynamic)

	MoveWithBallToPosition InstructionType = 4 // Move with the ball to a specific position (static)
	MoveWithBallToPlayer   InstructionType = 5 // Move with the ball to a specific player (dynamic)
	MoveWithBallToGoal     InstructionType = 6 // Move with the ball to some open space in the goal (dynamic)

	KickToPlayer   InstructionType = 7 // Kick the ball to a specific player (dynamic)
	KickToGoal     InstructionType = 8 // Kick the ball to some open space in the goal (dynamic)
	KickToPosition InstructionType = 9 // Kick the ball to a specific position (static)

	ReceiveBallFromPlayer InstructionType = 10 // Receive the ball from a specific player (dynamic). Before a kick is made, adjust the position to have a better chance of receiving the ball.
	ReceiveBallAtPosition InstructionType = 11 // Receive the ball at a expected position (dynamic). Make to be at the expected position when the ball is kicked. But adjust the position to have a better chance of receiving the ball after the kick.

	MeetBallFromPlayer InstructionType = 12 // Meet the ball (try to minimize the distance to the ball) from a specific player (dynamic)
	InterceptBall      InstructionType = 13 // Intercept the ball (try to minimize the distance to the ball) (dynamic)

	BlockEnemyPlayerFromPosition InstructionType = 14 // Body block an enemy player from a specific position (dynamic)
	BlockEnemyPlayerFromBall     InstructionType = 15 // Body block an enemy player from the ball (dynamic)
	BlockEnemyPlayerFromGoal     InstructionType = 16 // Body block an enemy player from the goal (dynamic). Make sure that the enemy does not have a clear shot at the goal.
	BlockEnemyPlayerFromPlayer   InstructionType = 17 // Body block an enemy player from a specific player (dynamic)
)

type Instruction struct {
	Type     InstructionType
	Id       ID
	Position Position
}

func (inst *Instruction) ToDTO() string {
	dto := fmt.Sprintf("Instruction{Id: %d, Position: %s}", inst.Id, inst.Position.ToDTO())
	return dto
}

type GamePlan struct {
	Valid        bool
	Team         Team
	Instructions []*Instruction
}

func NewGamePlan() *GamePlan {
	gp := GamePlan{}
	return &gp
}

func (gp *GamePlan) ToDTO() string {
	dto := fmt.Sprintf("GamePlan{Valid: %t, Team: %d, Instructions: [", gp.Valid, gp.Team)
	for _, instruction := range gp.Instructions {
		dto += instruction.ToDTO() + ", "
	}

	return dto
}
