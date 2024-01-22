package ai

import "fmt"

type RoleAssigner struct {
}

func NewRoleAssigner() *RoleAssigner {
	ra := &RoleAssigner{}
	return ra
}

type Roles struct {
}

func (ra *RoleAssigner) AssignRoles(plays *Plays) *Roles {
	fmt.Println("Hackspetten")
	return &Roles{}
}
