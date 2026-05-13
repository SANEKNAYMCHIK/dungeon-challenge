package dto

import (
	"dungeon-challenge/internal/domain"
	"fmt"
	"strconv"
	"strings"
)

func ToEventDomain(stringEvent string) (domain.Event, error) {
	data := strings.Split(stringEvent, " ")
	if len(data) < 3 {
		return domain.Event{}, fmt.Errorf("invalid event format")
	}
	TimeVal, err := ToTimeDomain(data[0])
	if err != nil {
		return domain.Event{}, err
	}
	IDVal, err := strconv.Atoi(data[1])
	if err != nil {
		return domain.Event{}, err
	}
	UserVal, err := strconv.Atoi(data[2])
	if err != nil {
		return domain.Event{}, err
	}
	return domain.Event{
		Time:  TimeVal,
		ID:    IDVal,
		User:  UserVal,
		Param: strings.Join(data[3:], " "),
	}, nil
}
