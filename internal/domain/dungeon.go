package domain

import "time"

type Dungeon struct {
	Levels   int
	Monsters int
	OpenAt   CustomTime
	Duration time.Duration
}
