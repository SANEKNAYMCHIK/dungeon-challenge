package domain

import "time"

type Dungeon struct {
	Floors   int
	Monsters int
	OpenAt   CustomTime
	CloseAt  CustomTime
	Duration time.Duration
}
