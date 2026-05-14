package domain

import "time"

type User struct {
	ID             int
	Health         UserHealth
	State          UserState
	FloorsTime     []time.Duration
	CurrentFloor   int
	MonstersKilled map[int]int
	ClearedFloor   map[int]bool
	FloorState     []bool
	FloorStartTime CustomTime
	BossKilled     bool
	BossStartTime  CustomTime
	BossDuration   time.Duration
	StartTime      CustomTime
	EndDuration    time.Duration
	Result         ReportHeader
}
