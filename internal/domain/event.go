package domain

type Event struct {
	Time  CustomTime
	ID    int
	User  int
	Param string
}
