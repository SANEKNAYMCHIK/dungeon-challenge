package domain

type EventType int

const (
	EventRegistered EventType = iota + 1
	EventInDungeon
	EventKilledMonster
	EventNextFloor
	EventPreviousFloor
	EventEnteredBossFloor
	EventKilledBoss
	EventLeftDungeon
	EventFailed
	EventGetHealth
	EventGetDamage
	EventDisqualified
	EventDead
	EventImpossibleMove
)
