package domain

import "time"

type Dungeon struct {
	Floors   int
	Monsters int
	OpenAt   CustomTime
	Duration time.Duration
}
