package ai

import (
	"fmt"
	"sync"
	"time"

	ai "github.com/LiU-SeeGoals/controller/internal/ai/activity"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type SlowBrain3 struct {
	SlowBrainComposition
	at_state int
	start    time.Time
	max_time time.Duration
}

func NewSlowBrain3(team info.Team) *SlowBrain3 {
	return &SlowBrain3{
		SlowBrainComposition: SlowBrainComposition{
			team: team,
		},
	}
}

func (m *SlowBrain3) Init(
	incoming <-chan info.GameInfo,
	activities *[]ai.Activity,
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
func (m *SlowBrain3) run() {

	for {
		// No need for slow brain to be fast
		time.Sleep(100 * time.Millisecond)

		//fmt.Println("slow, number of activities:", len(*m.activities))
			fmt.Println("done with action: ", m.team)
			m.AddActivity(ai.NewMoveToBall(m.team, 0))
			m.AddActivity(ai.NewMoveToBall(m.team, 1))
			m.AddActivity(ai.NewMoveToBall(m.team, 2))
			m.AddActivity(ai.NewMoveToBall(m.team, 3))
	}
}
