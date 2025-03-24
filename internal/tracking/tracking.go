package tracking

import (
	"github.com/LiU-SeeGoals/controller/internal/client"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type Tracking struct {
	robotTrackers       [2][info.TEAM_SIZE]*robotTracker
	ballTracker         *ballTracker
	possTracker         possessorTracker
}

func NewTracking(trackerAddress string) *Tracking {
	source := client.NewTrackerSource(trackerAddress)

	// Robot trackers
	var robotTrackers [2][info.TEAM_SIZE]*robotTracker
	for team := 0; team < 2; team++ {
		var id info.ID
		for id = 0; id < info.TEAM_SIZE; id++ {
			robotTrackers[team][id] = NewRobotTracker(source, info.Team(team), id)	
		}
	}

	t := &Tracking{
		ballTracker:         NewBallTracker(source),
		robotTrackers:       robotTrackers,
		possTracker:         possessorTracker{},
	}
	return t
}

func (tracking *Tracking) UpdateTracking(state *info.GameState) {

	// ##### ROBOTS #####
	var trackedRobots [2][info.TEAM_SIZE]info.RobotState
	for _, team := range []info.Team{info.Blue, info.Yellow} {
		var id info.ID
		for id = 0; id < info.TEAM_SIZE; id++ {

			// Track the robot
			trackedRobot := tracking.robotTrackers[team][id].GetTrackedRobot()
			trackedRobots[team][id] = trackedRobot
		}
	}
	state.UpdateRobots(trackedRobots)

	// ##### BALL #####
	trackedBall, _ := tracking.ballTracker.GetTrackedBall(state.BallPossessor)
	state.UpdateBall(trackedBall)

	// ##### BALL POSSESSOR #####
	possessor := tracking.possTracker.GetPossessor(state)
	state.UpdatePossessor(possessor)
}
