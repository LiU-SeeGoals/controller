package ai

import (
	"fmt"
	"sync"
	"time"

	ai "github.com/LiU-SeeGoals/controller/internal/ai/activity"
	"github.com/LiU-SeeGoals/controller/internal/helper"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type SlowBrain1 struct {
	SlowBrainComposition
	at_state int
	start    time.Time
	max_time time.Duration
}

func NewSlowBrain1(team info.Team) *SlowBrain1 {
	return &SlowBrain1{
		SlowBrainComposition: SlowBrainComposition{
			team: team,
		},
	}
}

func (m *SlowBrain1) Init(
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
func (m *SlowBrain1) run() {
	index := 0
	way_points := make([]info.Position, 0, 32)
	way_points = append(way_points, info.Position{X: 0, Y: 0})

	for {
		gi := <-m.incomingGameInfo

		ball := gi.State.Ball

		robot := gi.State.GetRobot(2, m.team)
		if robot == nil || !robot.IsActive() {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		pos, err := ball.GetEstimatedPosition()
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if index%2 == 0 {
			way_points = append(way_points, pos)
		} else {
			way_points = append(way_points, helper.Random_Position(gi))
			fmt.Println("Added waypoint: ", way_points[index])
			fmt.Println("Fieldsize: ", gi.FieldSize())
		}

		time.Sleep(100 * time.Millisecond)

		if len(way_points) == 0 {
			continue
		}
		if m.activities[2] == nil {
			fmt.Println("done with action: ", m.team)
			rPos, _ := robot.GetPosition()
			moveToPos := ai.NewMoveToPosition(m.team, robot.GetID(), way_points[index])
			action := moveToPos.GetMoveToAction(&gi)
			action.Dest.Angle = rPos.AngleToPosition(way_points[index])
			m.AddActivity(moveToPos)
			index = (index + 1) % len(way_points)
		}
	}
}
