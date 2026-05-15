package output

import (
	"dungeon-challenge/internal/domain"
	"os"
	"testing"
	"time"
)

func TestGetOutputLine(t *testing.T) {
	line := getOutputLine(
		domain.EventRegistered,
		domain.CustomTime{
			Time: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
		},
		1,
		"",
	)

	expected := "[10:00:00] Player [1] registered\n"

	if line != expected {
		t.Fatalf("expected %q got %q", expected, line)
	}
}

func TestWriterWrite(t *testing.T) {
	file, err := os.CreateTemp("", "output_test")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(file.Name())

	writer := &EventWriter{
		file: file,
	}

	writer.WriteEvent(
		domain.EventRegistered,
		domain.Event{
			ID: domain.EventRegistered,
			Time: domain.CustomTime{
				Time: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
			},
			User: 1,
		},
	)
}
