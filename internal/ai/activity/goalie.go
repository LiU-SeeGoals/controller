package ai

import (
	"fmt"

	"github.com/LiU-SeeGoals/controller/internal/action"
	. "github.com/LiU-SeeGoals/controller/internal/logger"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

// Goalie is the refactored keeper AI that follows a state machine (g.at_state)
// and tries to position itself based on the ball and a potential "shooter."
type Goalie struct {
	GenericComposition
	team     info.Team
	id       info.ID
	at_state int
}

func (g *Goalie) String() string {
	return fmt.Sprintf("Goalie(%d, %d)", g.team, g.id)
}

// NewGoalie creates a new Goalie struct.
func NewGoalie(team info.Team, id info.ID) *Goalie {
	return &Goalie{
		GenericComposition: GenericComposition{
			team: team,
			id:   id,
		},
		at_state: 0, // initial state
	}
}

// GetAction decides what the goalie should do each tick (frame), returning a single Action.
func (g *Goalie) GetAction(gi *info.GameInfo) action.Action {
	fmt.Println("Goalie")

	myTeam := gi.State.GetTeam(g.team)
	robot := myTeam[g.id]
	robotPos, err := robot.GetPosition()

	if err != nil {
		Logger.Errorf("Position retrieval failed - Robot: %v\n", err)
		return NewStop(g.id).GetAction(gi)
	}

	ballPos, err := gi.State.GetBall().GetPosition()

	if err != nil {
		Logger.Errorf("Ball position retrieval failed - Ball: %v\n", err)
		return NewStop(g.id).GetAction(gi)
	}


	// Prepare a MoveTo action for the goalie
	act := action.MoveTo{
		Id:   int(g.id),
		Team: g.team,
		Pos:  robotPos, // current position
	}

	// 1) Attempt to find the "shooter" from the opposing team
	shooter := g.findShooter(gi, ballPos)

	// 2) If there is a shooter, replicate the logic:
	//    - If ballPos.X <= 4000 => track shooter's Y
	//    - Else => move closer to the ball
	if shooter != nil {
		shooterPos, _ := shooter.GetPosition()
		if ballPos.X <= 4000 {
			if shooterPos.Y <= 0 {
				// Shooter is on the negative side (top in some coordinate systems)
				if shooterPos.Y <= -500 {
					act.Dest.Y = -350
				} else if shooterPos.Y <= -350 {
					act.Dest.Y = -250
				} else if shooterPos.Y <= -250 {
					act.Dest.Y = -150
				} else {
					act.Dest.Y = shooterPos.Y
				}
			} else {
				// Shooter is on the positive side (bottom in some coordinate systems)
				if shooterPos.Y >= 500 {
					act.Dest.Y = 350
				} else if shooterPos.Y >= 350 {
					act.Dest.Y = 250
				} else if shooterPos.Y >= 250 {
					act.Dest.Y = 150
				} else {
					act.Dest.Y = shooterPos.Y
				}
			}
			act.Dest.X = 4000
		} else {
			// The ball is "right" of X=4000, move closer to the ball
			act.Dest.X = ballPos.X + 25
			act.Dest.Y = ballPos.Y
		}
	} else {
		// 3) If NO shooter is found, do the "else" logic from the original code
		act.Dest.X = ballPos.X + 25
		act.Dest.Y = ballPos.Y
	}

	act.Dribble = false
	return &act
}

// Achieved returns whether this action is "complete".
// The goalie never really finishes, so we return false unless higher-level AI changes it.
func (g *Goalie) Achieved(*info.GameInfo) bool {
	return false
}

// findShooter checks the enemy team for any robot within a threshold distance of the ball.
// If found, returns that robot and true. Otherwise returns nil and false.
func (g *Goalie) findShooter(gi *info.GameInfo, ballPos info.Position) *info.Robot {
	// If your game only has two teams (0 and 1), you can identify the opposing team as:
	var enemyTeamID info.Team
	if g.team == 0 {
		enemyTeamID = 1
	} else {
		enemyTeamID = 0
	}

	enemyTeam := gi.State.GetTeam(enemyTeamID)
	const shooterThreshold = 500.0 // the "X distance" within which a robot is considered the shooter

	// Iterate over all robots in the enemy team
	for _, enemyRobot := range enemyTeam {
		enemyPos, err := enemyRobot.GetPosition()
		if err != nil {
			continue
		}
		dist := enemyPos.Distance(ballPos)
		if dist <= shooterThreshold {
			// Found a robot close enough to be considered "shooter"
			return enemyRobot
		}
	}
	// No enemy robot is close enough
	return nil
}

func (g *Goalie) GetID() info.ID {
	return g.id
}
