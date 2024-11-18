package constants

type MessageType int

const (
	COMMAND      = iota
	CALLBACK     = iota
	CONVERSATION = iota
)

func (m MessageType) String() string {
	return [...]string{"COMMAND", "CALLBACK", "CONVERSATION"}[m]
}

func (m MessageType) Int() int {
	return int(m)
}
