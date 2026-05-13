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
	Write(eventID domain.EventType, lineTime domain.CustomTime, userID int, params string) (int, error)
}

type DungeonRunner struct {
	dungeonInfo  domain.Dungeon
	eventsReader EventReader
	outputWriter EventWriter
	reportWriter *output.Writer
	users        map[int]*domain.User
}

func NewDungeonRunner(dungeonInfo domain.Dungeon, eventsReader EventReader, outputWriter *output.Writer, reportWriter *output.Writer) *DungeonRunner {
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

func (dr *DungeonRunner) getUser(id int) (*domain.User, bool) {
	user, exists := dr.users[id]
	return user, exists
}

func (dr *DungeonRunner) createNewUserWithDisqualification(event domain.Event) {
	dr.users[event.User] = &domain.User{
		ID:     event.User,
		Health: domain.MaxHealth,
		State:  domain.StateDisqualified,
		Result: domain.ReportHeaderDisqual,
	}
	dr.outputWriter.Write(domain.EventDisqualified, event.Time, event.User, "")
}

func (dr *DungeonRunner) HandleRegister(event domain.Event) {
	if _, exists := dr.getUser(event.User); !exists {
		dr.users[event.User] = &domain.User{}
		dr.users[event.User].ID = event.User
		dr.users[event.User].Health = domain.MaxHealth
		dr.users[event.User].State = domain.StateRegistered
		dr.outputWriter.Write(event.ID, event.Time, event.User, "")
	} else {
		dr.users[event.User].State = domain.StateDisqualified
		dr.users[event.User].Result = domain.ReportHeaderDisqual
		dr.outputWriter.Write(domain.EventDisqualified, event.Time, event.User, "")
	}
}

func (dr *DungeonRunner) HandleInDungeon(event domain.Event) {
	if user, exists := dr.getUser(event.User); exists {
		if user.State == domain.StateRegistered {
			user.State = domain.StateInDungeon
			dr.outputWriter.Write(event.ID, event.Time, event.User, "")
		} else {
			user.State = domain.StateDisqualified
			user.Result = domain.ReportHeaderDisqual
			dr.outputWriter.Write(domain.EventDisqualified, event.Time, event.User, "")
		}
	} else {
		dr.createNewUserWithDisqualification(event)
	}
}

func (dr *DungeonRunner) HandleHealth(event domain.Event) {
	if user, exists := dr.getUser(event.User); exists {
		if user.State == domain.StateInDungeon {
			healthVal, err := toInt(event.Param)
			if err != nil {
				log.Printf("error converting health value to int: %v", err)
				return
			}
			user.Health += domain.UserHealth(healthVal)
			if user.Health > domain.MaxHealth {
				user.Health = domain.MaxHealth
			}
			dr.outputWriter.Write(event.ID, event.Time, event.User, event.Param)
		} else {
			user.State = domain.StateDisqualified
			user.Result = domain.ReportHeaderDisqual
			dr.outputWriter.Write(domain.EventDisqualified, event.Time, event.User, "")
		}
	} else {
		dr.createNewUserWithDisqualification(event)
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
