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

func NewGoalie(team info.Team, id info.ID, dest info.Position) *MoveToPosition {
	return &MoveToPosition{
		team:            team,
		id:              id,
		target_position: dest,
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

	act.Dribble = false
	return &act
}
