package ai

import (
	"sync"
	"time"

	ai "github.com/LiU-SeeGoals/controller/internal/ai/activity"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type SlowBrainContainer struct {
	SlowBrainComposition
	at_state int
	start    time.Time
	max_time time.Duration
}

func NewSlowBrainContainer(team info.Team) *SlowBrainContainer {
	return &SlowBrainContainer{
		SlowBrainComposition: SlowBrainComposition{
			team: team,
		},
	}
}

func (m *SlowBrainContainer) Init(
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
func (m *SlowBrainContainer) run() {

	activityList := []ai.Activity{
		ai.NewMoveToPosition(m.team, 0, info.Position{X: 2000, Y: 0, Z: 0, Angle: 0}),
		ai.NewMoveToPosition(m.team, 0, info.Position{X: 0, Y: 2000, Z: 0, Angle: 0}),
		ai.NewMoveToPosition(m.team, 0, info.Position{X: 0, Y: 0, Z: 0, Angle: 0}),
	}

	queue := ai.NewActivityQueue(0, activityList)
	m.AddActivity(queue)
	// loop := ai.NewActivityLoop(0, activityList)
	// m.AddActivity(loop)

}
