package dto

import (
	"dungeon-challenge/internal/domain"
	"time"
)

func ToTimeDomain(stringTime string) (domain.CustomTime, error) {
	parsedTime, err := time.Parse(time.TimeOnly, stringTime[1:len(stringTime)-1])
	if err != nil {
		return domain.CustomTime{}, err
	}
	return domain.CustomTime{Time: parsedTime}, nil
}
