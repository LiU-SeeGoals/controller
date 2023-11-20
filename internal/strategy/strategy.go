package strategy

type Strategy interface {
	GetName() string
	GetCommand() string
}
