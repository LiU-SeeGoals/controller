package strategy

type Striker struct {
	name string
}

func NewStriker() *Striker {
	return &Striker{name: "striker"}
}

func (s Striker) GetName() string {
	return s.name
}

func (s Striker) GetCommand() string {
	// Implment code here
	return "strike the ball"
}
