package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type HandleReferee struct {
	team info.Team
}

func NewHandleReferee(team info.Team) *HandleReferee {
	return &HandleReferee{
		team: team,
	}
}

func (m *HandleReferee) GetActions(gi *info.GameInfo) []action.Action {
	switch gi.Status.GetGameEvent().RefCommand {
	case info.HALT:
		// Implement HALT logic
		return m.createStopActions(gi)
	case info.STOP:
		// Implement STOP logic
		return m.createStopActions(gi)
	case info.NORMAL_START:
		// Implement NORMAL_START logic
		return m.createStopActions(gi)
	case info.FORCE_START:
		// Implement FORCE_START logic
		return m.createStopActions(gi)
	case info.PREPARE_KICKOFF_YELLOW:
		// Implement PREPARE_KICKOFF_YELLOW logic
		return m.createStopActions(gi)
	case info.PREPARE_KICKOFF_BLUE:
		// Implement PREPARE_KICKOFF_BLUE logic
		return m.createStopActions(gi)
	case info.PREPARE_PENALTY_YELLOW:
		// Implement PREPARE_PENALTY_YELLOW logic
		return m.createStopActions(gi)
	case info.PREPARE_PENALTY_BLUE:
		// Implement PREPARE_PENALTY_BLUE logic
		return m.createStopActions(gi)
	case info.DIRECT_FREE_YELLOW:
		// Implement DIRECT_FREE_YELLOW logic
		return m.createStopActions(gi)
	case info.DIRECT_FREE_BLUE:
		// Implement DIRECT_FREE_BLUE logic
		return m.createStopActions(gi)
	case info.INDIRECT_FREE_YELLOW:
		// Implement INDIRECT_FREE_YELLOW logic
		return m.createStopActions(gi)
	case info.INDIRECT_FREE_BLUE:
		// Implement INDIRECT_FREE_BLUE logic
		return m.createStopActions(gi)
	case info.TIMEOUT_YELLOW:
		// Implement TIMEOUT_YELLOW logic
		return m.createStopActions(gi)
	case info.TIMEOUT_BLUE:
		// Implement TIMEOUT_BLUE logic
		return m.createStopActions(gi)
	case info.BALL_PLACEMENT_YELLOW:
		// Implement BALL_PLACEMENT_YELLOW logic
		return m.createStopActions(gi)
	case info.BALL_PLACEMENT_BLUE:
		// Implement BALL_PLACEMENT_BLUE logic
		return m.createStopActions(gi)
	default:
		return nil
	}
}

func (m *HandleReferee) createStopActions(gi *info.GameInfo) []action.Action {
	var actions []action.Action
	team := gi.State.GetTeam(m.team)
	for id := range team {
		actions = append(actions, &action.Stop{Id: int(id)})
	}
	return actions
}
