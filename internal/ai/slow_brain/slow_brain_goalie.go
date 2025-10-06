package ai

import (
	"fmt"
	"sync"
	"time"

	ai "github.com/LiU-SeeGoals/controller/internal/ai/activity"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type SlowBrainGoalie struct {
	SlowBrainComposition
	at_state int
	start    time.Time
	max_time time.Duration
}

func NewSlowBrainGoalie(team info.Team) *SlowBrainGoalie {
	return &SlowBrainGoalie{
		SlowBrainComposition: SlowBrainComposition{
			team: team,
		},
	}
}

func (m *SlowBrainGoalie) Init(
	incoming <-chan info.GameInfo,
	activities *[info.TEAM_SIZE]ai.Activity,
	lock *sync.Mutex,
	team info.Team,
) {
	m.incomingGameInfo = incoming
	m.activities = activities // store pointer directly
	m.activity_lock = lock
	m.team = team
	m.start = time.Now()

	go m.run()
}

// This is the main loop of the AI in this slow brain
func (m *SlowBrainGoalie) run() {
	var step = 0
	for {
		// No need for slow brain to be fast
		time.Sleep(100 * time.Millisecond)

		//fmt.Println("slow, number of activities:", len(*m.activities))

		if m.activities[0] == nil {
			fmt.Println("done with action: ", m.team)
			m.AddActivity(ai.NewGoalie(m.team, 0))
		}
		if m.activities[1] == nil {
			if step == 0 {
				fmt.Println("Adding move to ball for team: ", m.team)
				m.AddActivity(ai.NewMoveToBall(m.team, 1))
				step += 1
			} else if step == 1 {
				fmt.Println("Adding move with ball: ", m.team)
				m.AddActivity(ai.NewMoveWithBallToPosition(m.team, 1, info.Position{X: 250, Y: 0, Z: 0, Angle: 0}))
				step += 1
			} else if step == 2 {
				fmt.Println("Adding kick action to 1 for team: ", m.team)
				m.AddActivity(ai.NewKickTheBall(m.team, 1, info.Position{X: 2000, Y: 0}))
				step += 1
			}
		}
	}
}
