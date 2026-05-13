package usecase

import (
	"dungeon-challenge/internal/controller/output"
	"dungeon-challenge/internal/domain"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
)

type EventReader interface {
	ReadEvent() (domain.Event, error)
}

type EventWriter interface {
	WriteEvent(eventID domain.EventType, event domain.Event) (int, error)
	WriteImpossibleMove(eventID domain.EventType, event domain.Event, param string) (int, error)
}

type DungeonRunner struct {
	dungeonInfo  domain.Dungeon
	eventsReader EventReader
	outputWriter EventWriter
	reportWriter *output.EventWriter
	users        map[int]*domain.User
}

func NewDungeonRunner(dungeonInfo domain.Dungeon, eventsReader EventReader, outputWriter *output.EventWriter, reportWriter *output.EventWriter) *DungeonRunner {
	return &DungeonRunner{
		dungeonInfo:  dungeonInfo,
		eventsReader: eventsReader,
		outputWriter: outputWriter,
		reportWriter: reportWriter,
		users:        make(map[int]*domain.User),
	}
}

func toInt(param string) (int, error) {
	res, err := strconv.Atoi(param)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (dr *DungeonRunner) requireUser(id int) bool {
	_, exists := dr.users[id]
	return exists
}

func (dr *DungeonRunner) requireState(user *domain.User, expected domain.UserState) bool {
	if user.State != expected {
		return false
	}
	return true
}

func (dr *DungeonRunner) disqualifyUnregisteredUser(event domain.Event) {
	dr.users[event.User] = &domain.User{
		ID:     event.User,
		Health: domain.MaxHealth,
		State:  domain.StateDisqualified,
		Result: domain.ReportHeaderDisqual,
	}
	dr.outputWriter.WriteEvent(domain.EventDisqualified, event)
}

func (dr *DungeonRunner) HandleRegister(event domain.Event) {
	if !dr.requireUser(event.User) {
		dr.users[event.User] = &domain.User{
			ID:     event.User,
			Health: domain.MaxHealth,
			State:  domain.StateRegistered,
		}
		dr.outputWriter.WriteEvent(event.ID, event)
	} else {
		dr.users[event.User].State = domain.StateDisqualified
		dr.users[event.User].Result = domain.ReportHeaderDisqual
		dr.outputWriter.WriteEvent(domain.EventDisqualified, event)
	}
}

func (dr *DungeonRunner) HandleInDungeon(event domain.Event) {
	if dr.requireUser(event.User) {
		user := dr.users[event.User]
		if dr.requireState(user, domain.StateRegistered) {
			user.State = domain.StateInDungeon
			dr.outputWriter.WriteEvent(event.ID, event)
		} else {
			dr.outputWriter.WriteImpossibleMove(domain.EventImpossibleMove, event, event.ID.String())
			// user.State = domain.StateDisqualified
			// user.Result = domain.ReportHeaderDisqual
			// dr.outputWriter.Write(domain.EventDisqualified, event.Time, event.User, "")
		}
	} else {
		dr.disqualifyUnregisteredUser(event)
	}
}

func (dr *DungeonRunner) HandleHealth(event domain.Event) {
	if dr.requireUser(event.User) {
		user := dr.users[event.User]
		if dr.requireState(user, domain.StateInDungeon) {
			healthVal, err := toInt(event.Param)
			if err != nil {
				log.Printf("error converting health value to int: %v", err)
				return
			}
			user.Health += domain.UserHealth(healthVal)
			if user.Health > domain.MaxHealth {
				user.Health = domain.MaxHealth
			}
			dr.outputWriter.WriteEvent(event.ID, event)
		} else {
			dr.outputWriter.WriteImpossibleMove(domain.EventImpossibleMove, event, event.ID.String())
		}
	} else {
		dr.disqualifyUnregisteredUser(event)
	}
}

func (dr *DungeonRunner) executeEvent(event domain.Event) {
	switch event.ID {
	case domain.EventRegistered:
		dr.HandleRegister(event)
	case domain.EventInDungeon:
		dr.HandleInDungeon(event)
	case domain.EventKilledMonster:
		return
	case domain.EventNextFloor:
		return
	case domain.EventPreviousFloor:
		return
	case domain.EventEnteredBossFloor:
		return
	case domain.EventKilledBoss:
		return
	case domain.EventLeftDungeon:
		return
	case domain.EventFailed:
		return
	case domain.EventGetHealth:
		dr.HandleHealth(event)
	case domain.EventGetDamage:
		return
	case domain.EventDisqualified:
		return
	case domain.EventDead:
		return
	case domain.EventImpossibleMove:
		return
	default:
		return
	}
}

func (dr *DungeonRunner) Run() {
	for {
		event, err := dr.eventsReader.ReadEvent()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Printf("error reading event: %v", err)
			continue
		}
		fmt.Println(event)
		dr.executeEvent(event)
	}
}
