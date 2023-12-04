package strategy

type Defender struct {
	name string
}

func NewDefender() *Defender {
	return &Defender{name: "defender"}
}

func (s Defender) GetName() string {
	return s.name
}

func (s Defender) GetCommand() string {
	// Implment code here
	return "defend"
}
