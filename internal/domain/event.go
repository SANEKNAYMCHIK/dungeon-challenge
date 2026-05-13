package domain

type Event struct {
	Time  CustomTime
	ID    EventType
	User  int
	Param string
}
