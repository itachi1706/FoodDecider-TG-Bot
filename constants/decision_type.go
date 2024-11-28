package constants

type DecisionType int

const (
	GENERAL  = iota
	LOCATION = iota
	GROUP    = iota
	LOCATION_GROUP
)

func (m DecisionType) String() string {
	return [...]string{"GENERAL", "LOCATION", "GROUP", "LOCATION_GROUP"}[m]
}

func (m DecisionType) Int() int {
	return int(m)
}
