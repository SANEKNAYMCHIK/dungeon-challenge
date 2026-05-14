package domain

type User struct {
	ID             int
	Health         UserHealth
	State          UserState
	FloorsTime     []CustomDuration
	CurrentFloor   int
	MonstersKilled map[int]int
	ClearedFloor   map[int]bool
	FloorState     []bool
	FloorStartTime CustomTime
	BossKilled     bool
	BossStartTime  CustomTime
	BossDuration   CustomDuration
	StartTime      CustomTime
	EndDuration    CustomDuration
	Result         ReportHeader
}
