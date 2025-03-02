package ai

import (
	ai "github.com/LiU-SeeGoals/controller/internal/ai/activity"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type HandleReferee struct {
	team     info.Team
	prev_ref info.RefCommand
}

func NewHandleReferee(team info.Team) *HandleReferee {
	return &HandleReferee{
		team: team,
	}
}

// Link to rules:
// https://robocup-ssl.github.io/ssl-rules/sslrules.html#_referee_commands
func (m *HandleReferee) GetRefereeActivities(gi *info.GameInfo) []ai.Activity {
	switch gi.Status.GetGameEvent().RefCommand {
	case info.HALT:
		// When the halt command is issued, no robot is allowed to move or manipulate the ball.
		// There is a grace period of 2 seconds for the robots to brake.
		return m.createStopActivities(gi)
	case info.STOP:
		// When the stop command is issued, all robots have to slow down to less than 1.5 m/s.
		// Additionally, all robots have to keep at least 0.5 meters distance to the ball and
		// are not allowed to manipulate the ball.
		return m.createStopActivities(gi)
	case info.PREPARE_KICKOFF_YELLOW:
		// Implement PREPARE_KICKOFF_YELLOW logic
		return m.createStopActivities(gi)
	case info.PREPARE_KICKOFF_BLUE:
		// Implement PREPARE_KICKOFF_BLUE logic
		return m.createStopActivities(gi)
	case info.PREPARE_PENALTY_YELLOW:
		// Implement PREPARE_PENALTY_YELLOW logic
		return m.createStopActivities(gi)
	case info.PREPARE_PENALTY_BLUE:
		// Implement PREPARE_PENALTY_BLUE logic
		return m.createStopActivities(gi)
	case info.DIRECT_FREE_YELLOW:
		// Implement DIRECT_FREE_YELLOW logic
		return m.createStopActivities(gi)
	case info.DIRECT_FREE_BLUE:
		// Implement DIRECT_FREE_BLUE logic
		return m.createStopActivities(gi)
	case info.INDIRECT_FREE_YELLOW:
		// Implement INDIRECT_FREE_YELLOW logic
		return m.createStopActivities(gi)
	case info.INDIRECT_FREE_BLUE:
		// Implement INDIRECT_FREE_BLUE logic
		return m.createStopActivities(gi)
	case info.TIMEOUT_YELLOW:
		// Implement TIMEOUT_YELLOW logic
		return m.createStopActivities(gi)
	case info.TIMEOUT_BLUE:
		// Implement TIMEOUT_BLUE logic
		return m.createStopActivities(gi)
	case info.BALL_PLACEMENT_YELLOW:
		// Implement BALL_PLACEMENT_YELLOW logic
		return m.createStopActivities(gi)
	case info.BALL_PLACEMENT_BLUE:
		// Implement BALL_PLACEMENT_BLUE logic
		return m.createStopActivities(gi)
	default:
		return nil
	}
}

func (m *HandleReferee) createStopActivities(gi *info.GameInfo) []ai.Activity {
	var activities []ai.Activity
	team := gi.State.GetTeam(m.team)
	for id := range team {
		activities = append(activities, ai.NewStop(info.ID(id)))
	}
	return activities
}
