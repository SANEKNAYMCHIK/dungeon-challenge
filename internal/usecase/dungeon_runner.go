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
	WriteEvent(domain.EventType, domain.Event) (int, error)
	WriteImpossibleMove(domain.EventType, domain.Event, string) (int, error)
	WriteDeadUser(domain.EventType, domain.Event) (int, error)
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
			user.CurrentFloor++
			user.FloorStartTime = event.Time
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

func (dr *DungeonRunner) HandleDamage(event domain.Event) {
	if dr.requireUser(event.User) {
		user := dr.users[event.User]
		if dr.requireState(user, domain.StateInDungeon) {
			damageVal, err := toInt(event.Param)
			if err != nil {
				log.Printf("error converting damage value to int: %v", err)
				return
			}
			dr.outputWriter.WriteEvent(event.ID, event)
			user.Health -= domain.UserHealth(damageVal)
			if user.Health <= domain.MinHealth {
				user.Health = domain.MinHealth
				user.State = domain.StateDead
				user.Result = domain.ReportHeaderFail
				dr.outputWriter.WriteDeadUser(domain.EventDead, event)
			}
		} else {
			dr.outputWriter.WriteImpossibleMove(domain.EventImpossibleMove, event, event.ID.String())
		}
	} else {
		dr.disqualifyUnregisteredUser(event)
	}
}

func (dr *DungeonRunner) HandleFail(event domain.Event) {
	if dr.requireUser(event.User) {
		user := dr.users[event.User]
		user.State = domain.StateDisqualified
		user.Result = domain.ReportHeaderFail
		dr.outputWriter.WriteEvent(domain.EventFailed, event)
	} else {
		dr.disqualifyUnregisteredUser(event)
	}
}

func (dr *DungeonRunner) HandleKilling(event domain.Event) {
	if dr.requireUser(event.User) {
		user := dr.users[event.User]
		if dr.requireState(user, domain.StateInDungeon) {
			if user.MonstersKilled < dr.dungeonInfo.Monsters {
				user.MonstersKilled++
				dr.outputWriter.WriteEvent(domain.EventKilledMonster, event)
				if user.MonstersKilled == dr.dungeonInfo.Monsters {
					user.FloorsTime = append(user.FloorsTime, event.Time.Sub(user.FloorStartTime.Time))
				}
			} else {
				dr.outputWriter.WriteImpossibleMove(domain.EventImpossibleMove, event, event.ID.String())
			}
		} else {
			dr.outputWriter.WriteImpossibleMove(domain.EventImpossibleMove, event, event.ID.String())
		}
	} else {
		dr.disqualifyUnregisteredUser(event)
	}
}

func (dr *DungeonRunner) NextFloor(event domain.Event) {
	if dr.requireUser(event.User) {
		user := dr.users[event.User]
		if dr.requireState(user, domain.StateInDungeon) && user.MonstersKilled == dr.dungeonInfo.Monsters {
			user.CurrentFloor++
			user.MonstersKilled = 0
			user.FloorStartTime = event.Time
			dr.outputWriter.WriteEvent(domain.EventNextFloor, event)
		} else {
			dr.outputWriter.WriteImpossibleMove(domain.EventImpossibleMove, event, event.ID.String())
		}
	} else {
		dr.disqualifyUnregisteredUser(event)
	}
}

func (dr *DungeonRunner) PreviousFloor(event domain.Event) {
	if dr.requireUser(event.User) {
		user := dr.users[event.User]
		if dr.requireState(user, domain.StateInDungeon) && user.MonstersKilled == dr.dungeonInfo.Monsters {
			user.CurrentFloor++
			user.MonstersKilled = 0
			user.FloorStartTime = event.Time
			dr.outputWriter.WriteEvent(domain.EventNextFloor, event)
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
		dr.HandleKilling(event)
	case domain.EventNextFloor:
		dr.NextFloor(event)
	case domain.EventPreviousFloor:
		dr.PreviousFloor(event)
	case domain.EventEnteredBossFloor:
		return
	case domain.EventKilledBoss:
		return
	case domain.EventLeftDungeon:
		return
	case domain.EventFailed:
		dr.HandleFail(event)
	case domain.EventGetHealth:
		dr.HandleHealth(event)
	case domain.EventGetDamage:
		dr.HandleDamage(event)
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
