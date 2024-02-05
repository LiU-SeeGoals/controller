package ai

import "github.com/LiU-SeeGoals/controller/internal/gamestate"

type RoleAssigner struct {
}

func NewRoleAssigner() *RoleAssigner {
	ra := &RoleAssigner{}
	return ra
}

func (ra *RoleAssigner) AssignRoles(plays *[]Play, gamestate gamestate.GameState) *[]Role {

	var roles []Role
	// Loop through plays order by priority
	for _, play := range *plays {

		// Loop though roles needed for play and assign them to robots
		for _, role := range play.roles {
			// TODO logic for choosing robots to assign for roles
			role.Assign(role.AssignHeuristic(gamestate.GetTeam(gamestate.Yellow)))

		}

		if len(roles) == 6 {
			break
		}
	}

	return &roles
}
