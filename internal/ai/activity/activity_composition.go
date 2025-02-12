package ai

import (
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type ActivityComposition struct {
	team info.Team
	id   info.ID
}

// Here we have funciton that are common across multiple activities,
// such as calculating the distance between two points.
// or movement that is legal and not blocked by other players.
