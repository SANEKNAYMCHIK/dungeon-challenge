package usecase

import (
	"dungeon-challenge/internal/controller/output"
	"dungeon-challenge/internal/domain"
	"errors"
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

func (dr *DungeonRunner) requireState(user *domain.User, expected ...domain.UserState) bool {
	for i := range expected {
		if user.State == expected[i] {
			return true
		}
	}
	return false
}

func (dr *DungeonRunner) allFloorsAreCleared(floorsVals map[int]bool) bool {
	for i := 1; i < dr.dungeonInfo.Floors; i++ {
		if !floorsVals[i] {
			return false
		}
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
			ID:             event.User,
			Health:         domain.MaxHealth,
			State:          domain.StateRegistered,
			ClearedFloor:   make(map[int]bool),
			MonstersKilled: make(map[int]int),
			FloorState:     make([]bool, dr.dungeonInfo.Floors+1),
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
			user.StartTime = event.Time
			user.FloorState[user.CurrentFloor] = true
			dr.outputWriter.WriteEvent(event.ID, event)
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
		}
	} else {
		dr.disqualifyUnregisteredUser(event)
	}
}

func (dr *DungeonRunner) HandleFail(event domain.Event) {
	if dr.requireUser(event.User) {
		user := dr.users[event.User]
		user.State = domain.StateDisqualified
		user.Result = domain.ReportHeaderDisqual
		dr.outputWriter.WriteEvent(domain.EventFailed, event)
	} else {
		dr.disqualifyUnregisteredUser(event)
	}
}

func (dr *DungeonRunner) HandleKilling(event domain.Event) {
	if dr.requireUser(event.User) {
		user := dr.users[event.User]
		if dr.requireState(user, domain.StateInDungeon) {
			if !user.ClearedFloor[user.CurrentFloor] && user.CurrentFloor < dr.dungeonInfo.Floors {
				user.MonstersKilled[user.CurrentFloor]++
				dr.outputWriter.WriteEvent(domain.EventKilledMonster, event)
				if user.MonstersKilled[user.CurrentFloor] == dr.dungeonInfo.Monsters {
					user.ClearedFloor[user.CurrentFloor] = true
					user.FloorsTime = append(user.FloorsTime, event.Time.Sub(user.FloorStartTime.Time))
				}
			} else {
				dr.outputWriter.WriteImpossibleMove(domain.EventImpossibleMove, event, event.ID.String())
			}
		}
	} else {
		dr.disqualifyUnregisteredUser(event)
	}
}

func (dr *DungeonRunner) HandleNextFloor(event domain.Event) {
	if dr.requireUser(event.User) {
		user := dr.users[event.User]
		if dr.requireState(user, domain.StateInDungeon) {
			if user.ClearedFloor[user.CurrentFloor] {
				user.CurrentFloor++
				if !user.FloorState[user.CurrentFloor] {
					user.FloorStartTime = event.Time
				}
				dr.outputWriter.WriteEvent(domain.EventNextFloor, event)
			} else {
				dr.outputWriter.WriteImpossibleMove(domain.EventImpossibleMove, event, event.ID.String())
			}
		}
	} else {
		dr.disqualifyUnregisteredUser(event)
	}
}

func (dr *DungeonRunner) HandlePreviousFloor(event domain.Event) {
	if dr.requireUser(event.User) {
		user := dr.users[event.User]
		if dr.requireState(user, domain.StateInDungeon) {
			if user.CurrentFloor > 1 {
				user.CurrentFloor--
				dr.outputWriter.WriteEvent(domain.EventPreviousFloor, event)
			} else {
				dr.outputWriter.WriteImpossibleMove(domain.EventImpossibleMove, event, event.ID.String())
			}
		}
	} else {
		dr.disqualifyUnregisteredUser(event)
	}
}

func (dr *DungeonRunner) HandleEnteredBossFloor(event domain.Event) {
	if dr.requireUser(event.User) {
		user := dr.users[event.User]
		if dr.requireState(user, domain.StateInDungeon) {
			user.BossStartTime = event.Time
			dr.outputWriter.WriteEvent(domain.EventEnteredBossFloor, event)
		}
	} else {
		dr.disqualifyUnregisteredUser(event)
	}
}

func (dr *DungeonRunner) HandleKillingBoss(event domain.Event) {
	if dr.requireUser(event.User) {
		user := dr.users[event.User]
		if dr.requireState(user, domain.StateInDungeon) {
			if user.CurrentFloor == dr.dungeonInfo.Floors {
				user.BossDuration = event.Time.Sub(user.FloorStartTime.Time)
				user.BossKilled = true
				user.ClearedFloor[user.CurrentFloor] = true
				dr.outputWriter.WriteEvent(domain.EventKilledBoss, event)
			} else {
				dr.outputWriter.WriteImpossibleMove(domain.EventImpossibleMove, event, event.ID.String())
			}
		}
	} else {
		dr.disqualifyUnregisteredUser(event)
	}
}

func (dr *DungeonRunner) HandleLeftDungeon(event domain.Event) {
	if dr.requireUser(event.User) {
		user := dr.users[event.User]
		if dr.requireState(user, domain.StateInDungeon) {
			if dr.allFloorsAreCleared(user.ClearedFloor) && user.BossKilled {
				user.State = domain.StateFinished
				user.EndDuration = event.Time.Sub(user.StartTime.Time)
				dr.outputWriter.WriteEvent(domain.EventLeftDungeon, event)
			} else {
				dr.outputWriter.WriteImpossibleMove(domain.EventImpossibleMove, event, event.ID.String())
			}
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
		dr.HandleNextFloor(event)
	case domain.EventPreviousFloor:
		dr.HandlePreviousFloor(event)
	case domain.EventEnteredBossFloor:
		dr.HandleEnteredBossFloor(event)
	case domain.EventKilledBoss:
		dr.HandleKillingBoss(event)
	case domain.EventLeftDungeon:
		dr.HandleLeftDungeon(event)
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
		if !event.Time.Before(dr.dungeonInfo.CloseAt.Time) {
			if dr.requireUser(event.User) {
				user := dr.users[event.User]
				if !dr.requireState(user, domain.StateDead, domain.StateDisqualified, domain.StateFinished) {
					user.Result = domain.ReportHeaderFail
					user.State = domain.StateDead
				}
			} else {
				dr.disqualifyUnregisteredUser(event)
			}
		} else {
			dr.executeEvent(event)
		}
	}
}
