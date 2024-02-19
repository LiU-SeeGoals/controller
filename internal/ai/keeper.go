package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"gonum.org/v1/gonum/mat"
)

type Keeper struct {
	id        int
	gamestate *gamestate.GameState
}

func (ke *Keeper) AssignHeuristic(robot [gamestate.TEAM_SIZE]*gamestate.Robot) int {
	return 5
}

// NewKeeper creates a new instance of Keeper
func NewKeeper(gameState *gamestate.GameState) *Keeper {
	return &Keeper{
		id:        -1,
		gamestate: gameState,
	}
}

func (ke *Keeper) Assign(id int) {
	ke.id = id
}

func (ke *Keeper) NextStep() action.Action {
	act := &action.Move{}

	ball_pos := ke.gamestate.GetBall().GetPosition()
	goalie := ke.gamestate.GetRobot(ke.id, false)
	goaliePos := goalie.GetPosition()
	act.Pos = goaliePos
	act.Id = ke.id
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
