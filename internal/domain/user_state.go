package domain

type UserState int

const (
	StateRegistered UserState = iota
	StateInDungeon
	StateDead
	StateDisqualified
	StateFinished
)
