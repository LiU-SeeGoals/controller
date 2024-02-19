package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"gonum.org/v1/gonum/mat"
)

type Baller struct {
	id        int
	gamestate *gamestate.GameState
	condition *mat.VecDense
}

// NewBaller creates a new instance of Baller
func NewBaller(gameState *gamestate.GameState, condition *mat.VecDense) *Baller {
	return &Baller{
		id:        -1,
		condition: condition,
		gamestate: gameState,
	}
}

// TODO complete this function
func (ba *Baller) AssignHeuristic(robot [gamestate.TEAM_SIZE]*gamestate.Robot) int {
	return 1
}

func (ba *Baller) Assign(id int) {
	ba.id = id
}

func (ba *Baller) NextStep() action.Action {
	act := &action.Move{}

	ball_pos := ba.gamestate.GetBall().GetPosition()
	goalie := ba.gamestate.GetRobot(ba.id, false)
	goaliePos := goalie.GetPosition()
	act.Pos = goaliePos
	act.Id = ba.id
	halfGoalHeight := 1000

	dest := mat.NewVecDense(3, nil)

	if ball_pos.AtVec(0) < float64(3500) {
		if ball_pos.AtVec(1) > -float64(halfGoalHeight) && ball_pos.AtVec(1) < float64(halfGoalHeight) {
			dest.SetVec(0, 3500)
			dest.SetVec(1, ball_pos.AtVec(1))
		} else if ball_pos.AtVec(1) < 0 {
			dest.SetVec(0, 3500)
			dest.SetVec(1, -1000)
		} else {
			dest.SetVec(0, 3500)
			dest.SetVec(1, 1000)
		}
	} else if ball_pos.AtVec(1) < 0 {
		dest.SetVec(0, ball_pos.AtVec(0))
		dest.SetVec(1, -1000)
	} else {
		dest.SetVec(0, ball_pos.AtVec(0))
		dest.SetVec(1, 1000)
	}

	act.Dest = dest
	act.Dribble = true
	return act
}
