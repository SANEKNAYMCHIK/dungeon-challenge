package dto

import (
	"dungeon-challenge/internal/domain"
	"time"
)

type Dungeon struct {
	Floors   int               `json:"Floors"`
	Monsters int               `json:"Monsters"`
	OpenAt   domain.CustomTime `json:"OpenAt"`
	Duration int               `json:"Duration"`
}

func (d *Dungeon) ToDomain() domain.Dungeon {
	return domain.Dungeon{
		Floors:   d.Floors,
		Monsters: d.Monsters,
		OpenAt:   d.OpenAt,
		Duration: time.Duration(d.Duration) * time.Hour,
	}
}
