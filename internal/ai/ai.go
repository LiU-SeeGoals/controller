package ai

import (
	"sync"

	"github.com/LiU-SeeGoals/controller/internal/action"
	ai "github.com/LiU-SeeGoals/controller/internal/ai/activity"
	"github.com/LiU-SeeGoals/controller/internal/helper"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type SlowBrain interface {
	Init(incoming <-chan info.GameInfo, activities *[]ai.Activity, lock *sync.Mutex, team info.Team)
}

type FastBrain interface {
	Init(incoming <-chan info.GameInfo,
		activities *[]ai.Activity,
		lock *sync.Mutex,
		outgoing chan<- []action.Action,
		team info.Team,
	)
}

type Ai struct {
	team             info.Team
	slow_brain       SlowBrain
	fast_brain       FastBrain
	gameInfoSenderSB chan<- info.GameInfo
	gameInfoSenderFB chan<- info.GameInfo
	actionReceiver   chan []action.Action
	activities       *[]ai.Activity // Shared slice of Activity
	activity_lock    *sync.Mutex    // Shared mutex for synchronization
}

// Constructor for the AI
func NewAi(team info.Team, slowBrain SlowBrain, fastBrain FastBrain) *Ai {
	// Create a shared slice of Activity and a mutex
	activities := &[]ai.Activity{}
	lock := &sync.Mutex{}

	gameInfoSenderSB, gameInfoReceiverSB := helper.NB_KeepLatestChan[info.GameInfo]()
	gameInfoSenderFB, gameInfoReceiverFB := helper.NB_KeepLatestChan[info.GameInfo]()
	actionReceiver := make(chan []action.Action)

	// Initialize SlowBrain and FastBrain with the shared resources
	slowBrain.Init(gameInfoReceiverSB, activities, lock, team)
	fastBrain.Init(gameInfoReceiverFB, activities, lock, actionReceiver, team)

	// Construct the AI object
	ai := &Ai{
		team:             team,
		slow_brain:       slowBrain,
		fast_brain:       fastBrain,
		activities:       activities,
		gameInfoSenderSB: gameInfoSenderSB,
		gameInfoSenderFB: gameInfoSenderFB,
		activity_lock:    lock,
		actionReceiver:   actionReceiver,
	}
	return ai
}

// Decides on new actions for the robots
func (ai *Ai) GetActions(gi *info.GameInfo) []action.Action {

	// Send the game state copy to the slow brain
	ai.gameInfoSenderSB <- *gi

	// Send the game state to the fast brain
	ai.gameInfoSenderFB <- *gi

	// Get the actions from the fast brain, this will block until the fast brain has decided on actions
	actions := <-ai.actionReceiver
	if len(actions) > 0 {
		//fmt.Println(actions[0])
	}
	//fmt.Println(actions[0])
	return actions
}
