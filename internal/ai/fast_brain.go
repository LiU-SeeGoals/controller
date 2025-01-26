package ai

import (
	"sync"

	"github.com/LiU-SeeGoals/controller/internal/action"
	ai "github.com/LiU-SeeGoals/controller/internal/ai/activity"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type FastBrainGO struct {
	team             info.Team
	incomingGameInfo <-chan info.GameInfo
	outgoingActions  chan<- []action.Action
	activities       []ai.Activity
	activity_lock    *sync.Mutex // Shared mutex for synchronization
}

func NewFastBrainGO() *FastBrainGO {
	return &FastBrainGO{}
}

func (fb *FastBrainGO) Init(incoming <-chan info.GameInfo,
	activities *[]ai.Activity,
	lock *sync.Mutex,
	outgoing chan<- []action.Action,
	team info.Team,
) {

	fb.incomingGameInfo = incoming
	fb.outgoingActions = outgoing

	go fb.Run()
}

func (fb *FastBrainGO) Run() {
	gameInfo := info.GameInfo{}
	for {
		var activitiesCopy []ai.Activity
		var actions []action.Action

		for i := range activitiesCopy {

			// If the activity is done, remove it
			if fb.activities[i].Achieved(&gameInfo) {
				fb.activity_lock.Lock()
				fb.activities = append(fb.activities[:i], fb.activities[i+1:]...)
				fb.activity_lock.Unlock()
			} else {
				actions = append(actions, fb.activities[i].GetAction(&gameInfo))
			}

		}

		// Send the actions to the AI
		fb.outgoingActions <- actions
		// fmt.Println("FastBrainGO: Sent actions")

	}
}
