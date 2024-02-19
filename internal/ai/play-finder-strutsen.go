package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/gamestate"
	"gonum.org/v1/gonum/mat"
)

type PlayFinder struct {
}

func NewPlayFinder() *PlayFinder {
	pf := &PlayFinder{}
	return pf
}

// struct to be extended, max/min robots etc..
type Play struct {
	roles []Role
}

func (pf *PlayFinder) FindPlays(gameAnalysis *GameAnalysis, gamestateObj *gamestate.GameState) []Play {
	var plays []Play

	// Keeper play
	if true {
		var roles []Role

		roles = append(roles, NewKeeper(gamestateObj))
		play := Play{
			roles: roles,
		}
		plays = append(plays, play)
	}

	// move ball
	if true {
		vec := mat.NewVecDense(2, nil)
		var roles []Role
		roles = append(roles, NewBaller(gamestateObj, vec))
		play := Play{
			roles: roles,
		}
		plays = append(plays, play)
	}

	return plays
}
