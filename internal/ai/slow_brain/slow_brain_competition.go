package ai

import (
	"fmt"
	"sync"
	"time"

	ai "github.com/LiU-SeeGoals/controller/internal/ai/activity"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

// ========================================================
// SlowBrainCompetition is a slow brain used for competition
// ========================================================

type SlowBrainCompetition struct {
	SlowBrainComposition
	HandleReferee

	at_state int
	start    time.Time
	max_time time.Duration
	team     info.Team
	prev_ref info.RefCommand
}

func NewSlowBrainCompetition(team info.Team) *SlowBrainCompetition {
	return &SlowBrainCompetition{
		SlowBrainComposition: SlowBrainComposition{
			team: team,
		},
		HandleReferee: HandleReferee{
			team: team,
		},
		team: team,
	}
}

func (m *SlowBrainCompetition) Init(
	incoming <-chan info.GameInfo,
	activities *[info.TEAM_SIZE]ai.Activity,
	lock *sync.Mutex,
	team info.Team,
) {
	m.incomingGameInfo = incoming
	m.activities = activities // store pointer directly
	m.activity_lock = lock
	m.start = time.Now()

	go m.run()
}

// This is the main loop of the AI in this slow brain
func (m *SlowBrainCompetition) run() {
	way_points := []info.Position{
		{X: 0, Y: 4400, Z: 0, Angle: 0},
		{X: 0, Y: -4400, Z: 0, Angle: 0},
	}

	enemy_goal := 0
	fmt.Println("SlowBrainCompetition: starting")

	for {
		// No need for slow brain to be fast
		time.Sleep(100 * time.Millisecond)

		gameInfo := <-m.incomingGameInfo

		if gameInfo.Status.GetBlueTeamOnPositiveHalf() && m.team == info.Blue {
			enemy_goal = 1
		} else {
			enemy_goal = 0

		}

		referee_activities := m.GetRefereeActivities(&gameInfo)

		for _, act := range referee_activities {
			if act != nil { // Check for a non-nil activity
				fmt.Println("adding referee activity")
				m.ReplaceActivities(referee_activities)
				m.at_state = REFEREE
				continue
			}
		}

		// If we are EXITING the REFEREE state, we need to clear the activities
		if m.at_state == REFEREE {
			fmt.Println("clearing activities")
			m.ClearActivities()
			m.at_state = RUNNING
		}

		// Set robot to goalie
		if m.activities[0] == nil {
			fmt.Println("done with action: ", m.team)
			m.AddActivity(ai.NewGoalie(m.team, 0))
		}

		// The other robot is doing all the work

		// The logic for the other robot
		// 1. Chaise the ball
		// 2. If get the ball, dribble to a position in front of thier goal
		// 3. Kick the ball to the goal
		// 4. Repeat

		if m.activities[1] == nil {
			fmt.Println("done with action: ", m.team)

			// If we have the ball, then dribble to the enemy goal
			if gameInfo.State.GetBall().GetPossessor().GetID() == 1 {
				m.AddActivity(ai.NewMoveWithBallToPosition(m.team, 1, way_points[enemy_goal]))

			} else {
				m.AddActivity(ai.NewMoveToBall(m.team, 1))
			}

		}
	}
}
