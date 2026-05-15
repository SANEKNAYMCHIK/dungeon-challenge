package usecase

import (
	"dungeon-challenge/internal/domain"
	"testing"
	"time"
)

func makeRunner() *DungeonRunner {
	return &DungeonRunner{
		dungeonInfo: domain.Dungeon{
			Floors:   3,
			Monsters: 2,
		},
		users:        make(map[int]*domain.User),
		outputWriter: &MockEventWriter{},
		reportWriter: &MockReportWriter{},
	}
}

func makeTime() domain.CustomTime {
	return domain.CustomTime{
		Time: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
	}
}

func TestHandleRegister(t *testing.T) {
	dr := makeRunner()

	event := domain.Event{
		ID:   domain.EventRegistered,
		User: 1,
	}

	dr.HandleRegister(event)

	user, exists := dr.users[1]

	if !exists {
		t.Fatal("user not created")
	}

	if user.State != domain.StateRegistered {
		t.Fatalf("expected registered state")
	}
}

func TestHandleInDungeon(t *testing.T) {
	dr := makeRunner()

	dr.users[1] = &domain.User{
		ID:         1,
		State:      domain.StateRegistered,
		FloorState: make([]bool, dr.dungeonInfo.Floors+1),
	}

	event := domain.Event{
		ID:   domain.EventInDungeon,
		User: 1,
		Time: makeTime(),
	}

	dr.HandleInDungeon(event)

	user := dr.users[1]

	if user.State != domain.StateInDungeon {
		t.Fatalf("expected in dungeon")
	}

	if user.CurrentFloor != 1 {
		t.Fatalf("expected floor 1")
	}
}

func TestHandleHealth(t *testing.T) {
	dr := makeRunner()

	dr.users[1] = &domain.User{
		ID:     1,
		State:  domain.StateInDungeon,
		Health: 50,
	}

	event := domain.Event{
		ID:    domain.EventGetHealth,
		User:  1,
		Param: "20",
	}

	dr.HandleHealth(event)

	if dr.users[1].Health != 70 {
		t.Fatalf("expected 70 hp")
	}
}

func TestHandleDamage(t *testing.T) {
	dr := makeRunner()

	dr.users[1] = &domain.User{
		ID:     1,
		State:  domain.StateInDungeon,
		Health: 100,
	}

	event := domain.Event{
		ID:    domain.EventGetDamage,
		User:  1,
		Param: "30",
	}

	dr.HandleDamage(event)

	if dr.users[1].Health != 70 {
		t.Fatalf("expected 70 hp")
	}
}

func TestHandleDeath(t *testing.T) {
	dr := makeRunner()

	dr.users[1] = &domain.User{
		ID:     1,
		State:  domain.StateInDungeon,
		Health: 10,
	}

	event := domain.Event{
		ID:    domain.EventGetDamage,
		User:  1,
		Param: "20",
	}

	dr.HandleDamage(event)

	user := dr.users[1]

	if user.State != domain.StateDead {
		t.Fatalf("expected dead state")
	}
}

func TestHandleKillMonster(t *testing.T) {
	dr := makeRunner()

	dr.users[1] = &domain.User{
		ID:             1,
		State:          domain.StateInDungeon,
		CurrentFloor:   1,
		ClearedFloor:   make(map[int]bool),
		MonstersKilled: make(map[int]int),
	}

	event := domain.Event{
		ID:   domain.EventKilledMonster,
		User: 1,
		Time: makeTime(),
	}

	dr.HandleKilling(event)

	if dr.users[1].MonstersKilled[1] != 1 {
		t.Fatalf("monster not counted")
	}
}

func TestHandleNextFloor(t *testing.T) {
	dr := makeRunner()

	dr.users[1] = &domain.User{
		ID:           1,
		State:        domain.StateInDungeon,
		CurrentFloor: 1,
		ClearedFloor: map[int]bool{
			1: true,
		},
		FloorState: make([]bool, 10),
	}

	event := domain.Event{
		ID:   domain.EventNextFloor,
		User: 1,
		Time: makeTime(),
	}

	dr.HandleNextFloor(event)

	if dr.users[1].CurrentFloor != 2 {
		t.Fatalf("expected floor 2")
	}
}

func TestHandlePreviousFloor(t *testing.T) {
	dr := makeRunner()

	dr.users[1] = &domain.User{
		ID:           1,
		State:        domain.StateInDungeon,
		CurrentFloor: 3,
	}

	event := domain.Event{
		ID:   domain.EventPreviousFloor,
		User: 1,
	}

	dr.HandlePreviousFloor(event)

	if dr.users[1].CurrentFloor != 2 {
		t.Fatalf("expected floor 2")
	}
}

func TestHandleBossKill(t *testing.T) {
	dr := makeRunner()

	dr.users[1] = &domain.User{
		ID:           1,
		State:        domain.StateInDungeon,
		CurrentFloor: 3,
		ClearedFloor: make(map[int]bool),
	}

	event := domain.Event{
		ID:   domain.EventKilledBoss,
		User: 1,
		Time: makeTime(),
	}

	dr.HandleKillingBoss(event)

	if !dr.users[1].BossKilled {
		t.Fatalf("boss should be killed")
	}
}

func TestHandleLeftDungeon(t *testing.T) {
	dr := makeRunner()

	dr.users[1] = &domain.User{
		ID:           1,
		State:        domain.StateInDungeon,
		CurrentFloor: 3,
		BossKilled:   true,
		ClearedFloor: map[int]bool{
			1: true,
			2: true,
			3: true,
		},
		StartTime: makeTime(),
	}

	event := domain.Event{
		ID:   domain.EventLeftDungeon,
		User: 1,
		Time: domain.CustomTime{
			Time: makeTime().Add(time.Hour),
		},
	}

	dr.HandleLeftDungeon(event)

	if dr.users[1].State != domain.StateFinished {
		t.Fatalf("expected finished state")
	}
}

func TestImpossibleMove(t *testing.T) {
	dr := makeRunner()

	dr.users[1] = &domain.User{
		ID:    1,
		State: domain.StateRegistered,
	}

	event := domain.Event{
		ID:   domain.EventNextFloor,
		User: 1,
	}

	dr.HandleNextFloor(event)

	if dr.users[1].CurrentFloor != 0 {
		t.Fatalf("floor should not change")
	}
}

func TestDungeonClosed(t *testing.T) {
	dr := makeRunner()

	dr.dungeonInfo.CloseAt = domain.CustomTime{
		Time: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
	}

	dr.users[1] = &domain.User{
		ID:    1,
		State: domain.StateInDungeon,
	}

	event := domain.Event{
		Time: domain.CustomTime{
			Time: time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC),
		},
	}

	if !event.Time.Before(dr.dungeonInfo.CloseAt.Time) {
		user := dr.users[1]
		user.State = domain.StateDead
	}

	if dr.users[1].State != domain.StateDead {
		t.Fatalf("user should be dead")
	}
}
