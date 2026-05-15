package dto

import (
	"dungeon-challenge/internal/domain"
	"testing"
)

func TestToEventDomain(t *testing.T) {
	event, err := ToEventDomain("[10:00:00] 42 1")

	if err != nil {
		t.Fatal(err)
	}

	if event.User != 42 {
		t.Fatalf("expected user 42")
	}

	if event.ID != domain.EventRegistered {
		t.Fatalf("wrong event id")
	}
}

func TestInvalidEvent(t *testing.T) {
	_, err := ToEventDomain("bad input")

	if err == nil {
		t.Fatalf("expected error")
	}
}
