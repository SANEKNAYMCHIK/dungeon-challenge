package domain

import "time"

type User struct {
	ID             int
	Health         UserHealth
	State          UserState
	FloorsTime     []time.Duration
	CurrentFloor   int
	FloorStartTime CustomTime
	BossKilled     bool
	BossStartTime  CustomTime
	BossEndTime    CustomTime
	StartTime      CustomTime
	EndTime        CustomTime
	Result         ReportHeader
}
