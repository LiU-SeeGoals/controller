package ai

import (
	"fmt"

	"github.com/LiU-SeeGoals/controller/internal/action"
	"github.com/LiU-SeeGoals/controller/internal/info"
)

type Goalie struct {
	team            info.Team
	id              info.ID
	target_position info.Position
}

func NewGoalie(team info.Team, id info.ID) *MoveToPosition {
	return &MoveToPosition{
		team: team,
		id:   id,
	}
}

func (fb *Goalie) GetAction(inst *info.Instruction, gs *info.GameState) action.Action {
	// todo: add collision avoidance
	myTeam := gs.GetTeam(fb.team)
	robot := myTeam[inst.Id]
	shooter := myTeam[1]

	robotPos := robot.GetPosition()
	shooterPos := shooter.GetPosition()
	if !robot.IsActive() {
		return nil
	}
	act := action.MoveTo{}
	act.Id = int(robot.GetID())
	act.Team = fb.team
	act.Pos = robotPos

	ballPos, _ := gs.GetBall().GetPositionTime()

	//dright := math.Abs(float64(shooterPos.Y - 500))

	//drightgoalie := math.Abs(float64(shooterPos.Y - robotPos.Y))

	//angleMiddle := math.Atan(float64(math.abs(shooterPos.Y)/2000))

	//dleft := math.Abs(float64(shooterPos.Y + 500))
	//dleftgoalie := math.Abs(float64(shooterPos.Y + i))

	//angleGoalie := math.Atan(float64(dleftgoalie/2000))
	//anglePose := math.Atan(float64(dleft/2000))

	fmt.Println("BallPos", ballPos)
	if ballPos.X <= 4000 {
		if shooterPos.Y <= 0 {
			if shooterPos.Y <= -500 {
				act.Dest.Y = -350
			} else if shooterPos.Y <= -350 {
				act.Dest.Y = -250
			} else if shooterPos.Y <= -250 {
				act.Dest.Y = -150
			} else {
				act.Dest.Y = shooterPos.Y
			}

		} else {
			if shooterPos.Y >= 500 {
				act.Dest.Y = 350
			} else if shooterPos.Y >= 350 {
				act.Dest.Y = 250
			} else if shooterPos.Y >= 250 {
				act.Dest.Y = 150
			} else {
				act.Dest.Y = shooterPos.Y
			}
		}

		act.Dest.X = 4000
	} else {

		act.Dest.Y = ballPos.Y
		act.Dest.X = ballPos.X + 25
	}

	// // next state
	// if g.at_state == 0 {
	// 	return []*info.Instruction{
	// 		{Type: info.MoveToPosition, Id: 0, Position: info.Position{X: 4000, Y: 0}},
	// 		{Type: info.MoveToBall, Id: 1},
	// 	}
	// } else if g.at_state == 1 {
	// 	return []*info.Instruction{
	// 		{Type: info.MoveWithBallToPosition, Id: 1, Position: info.Position{X: 2500, Y: 0}},
	// 	}
	// } else if g.at_state == 2 {

	// 	return []*info.Instruction{

	// 		{Type: info.MoveWithBallToPosition, Id: 1, Position: info.Position{X: 3550, Y: 1000}},
	// 		{Type: info.Goalie, Id: 0},
	// 	}
	// } else if g.at_state == 3 {

	// 	return []*info.Instruction{

	// 		{Type: info.MoveWithBallToPosition, Id: 1, Position: info.Position{X: 3500, Y: -1000}},
	// 		{Type: info.Goalie, Id: 0},
	// 	}
	// } else {
	// 	return []*info.Instruction{}
	// }

	// // Check if current state is done
	// robot_pos := gs.GetRobot(info.ID(1), g.team).GetPosition()
	// ball_pos := gs.GetBall().GetPosition()

	// dxBall := float64(robot_pos.X - ball_pos.X)
	// dyBall := float64(robot_pos.Y - ball_pos.Y)
	// distanceBall := math.Sqrt(math.Pow(dxBall, 2) + math.Pow(dyBall, 2))

	// dxPos := float64(robot_pos.X - (2500))
	// dyPos := float64(robot_pos.Y)
	// distancePos := math.Sqrt(math.Pow(dxPos, 2) + math.Pow(dyPos, 2))
	// fmt.Println("Gamestate:")
	// fmt.Println(g.at_state)
	// if g.at_state == 0 {
	// 	if distanceBall < 1 {
	// 		g.at_state = 1
	// 	}
	// } else if g.at_state == 1 {

	// 	if distancePos < 100 {

	// 		g.at_state = 2
	// 	}
	// } else if g.at_state == 2 {
	// 	if dyPos > 950 {

	// 		g.at_state = 3
	// 	}
	// } else if g.at_state == 3 {

	// 	if dyPos < -950 {

	// 		g.at_state = 2
	// 	}
	// }

	act.Dribble = false
	return &act
}

func (g *Goalie) Archived(gs *info.GameState) bool {
	// The goalie is never done with its action
	// This means that only slow brain can change the action
	return false
}
