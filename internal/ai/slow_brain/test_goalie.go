
type Goalie struct {
	team     info.Team
	at_state int
	start    time.Time
	max_time time.Duration
}

func NewGoalie(team info.Team) *Goalie {
	return &Goalie{
		team:     team,
		max_time: 30 * time.Second,
		at_state: -1,
	}
}

func (g *Goalie) Run() []*info.Instruction {
	if g.at_state == -1 {
		g.start = time.Now()
		g.at_state = 0

	}

	if g.at_state == 0 {
		return []*info.Instruction{
			{Type: info.MoveToPosition, Id: 0, Position: info.Position{X: 4000, Y: 0}},
			{Type: info.MoveToBall, Id: 1},
		}
	} else if g.at_state == 1 {
		return []*info.Instruction{
			{Type: info.MoveWithBallToPosition, Id: 1, Position: info.Position{X: 2500, Y: 0}},
		}
	} else if g.at_state == 2 {

		return []*info.Instruction{

			{Type: info.MoveWithBallToPosition, Id: 1, Position: info.Position{X: 3550, Y: 1000}},
			{Type: info.Goalie, Id: 0},
		}
	} else if g.at_state == 3 {

		return []*info.Instruction{

			{Type: info.MoveWithBallToPosition, Id: 1, Position: info.Position{X: 3500, Y: -1000}},
			{Type: info.Goalie, Id: 0},
		}
	} else {
		return []*info.Instruction{}
	}

}

func (g *Goalie) Archived(gs *info.GameState) int {
	robot_pos := gs.GetRobot(info.ID(1), g.team).GetPosition()
	ball_pos := gs.GetBall().GetPosition()

	dxBall := float64(robot_pos.X - ball_pos.X)
	dyBall := float64(robot_pos.Y - ball_pos.Y)
	distanceBall := math.Sqrt(math.Pow(dxBall, 2) + math.Pow(dyBall, 2))

	dxPos := float64(robot_pos.X - (2500))
	dyPos := float64(robot_pos.Y)
	distancePos := math.Sqrt(math.Pow(dxPos, 2) + math.Pow(dyPos, 2))
	fmt.Println("Gamestate:")
	fmt.Println(g.at_state)
	if g.at_state == 0 {
		if distanceBall < 1 {
			g.at_state = 1
		}
	} else if g.at_state == 1 {

		if distancePos < 100 {

			g.at_state = 2
		}
	} else if g.at_state == 2 {
		if dyPos > 950 {

			g.at_state = 3
		}
	} else if g.at_state == 3 {

		if dyPos < -950 {

			g.at_state = 2
		}
	}

	return RUNNING
}
