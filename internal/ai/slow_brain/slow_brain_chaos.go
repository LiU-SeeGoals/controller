package ai

import (
	"fmt"
	"sync"
	"time"

	ai "github.com/LiU-SeeGoals/controller/internal/ai/activity"
	"github.com/LiU-SeeGoals/controller/internal/helper"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type SlowBrainChaos struct {
	SlowBrainComposition
	at_state int
	start    time.Time
	max_time time.Duration
}

func NewSlowBrainChaos(team info.Team) *SlowBrainChaos {
	return &SlowBrainChaos{
		SlowBrainComposition: SlowBrainComposition{
			team: team,
		},
	}
}

func (m *SlowBrainChaos) Init(
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

	fmt.Println("SBC Initiated")

	go m.run()
}

// This is the main loop of the AI in this slow brain
func (m *SlowBrainChaos) run() {
	for {
		gi := <-m.incomingGameInfo
		team := gi.State.GetTeam(m.team)

		for i := 0; i < int(info.TEAM_SIZE); i++ {
			robot := team[i]

			// Only assign a new random target when the previous activity has completed
			// and FastBrain has cleared the slot (activities[i] == nil)
			// m.activites[2]
			if m.activities[i] == nil {
				dest := helper.Random_Position(gi)
				pos, _ := robot.GetPosition()
				dest.Angle = pos.AngleToPosition(dest)
				fmt.Println("Assign random target:", robot.GetID(), "->", dest.String())
				m.AddActivity(ai.NewMoveToPosition(robot.GetTeam(), robot.GetID(), dest))
			}
		}
	}
}
