package report

import (
	"dungeon-challenge/internal/domain"
	"os"
	"testing"
	"time"
)

func TestGetOutputLine(t *testing.T) {
	line := getOutputLine(
		domain.ReportHeaderSuccess,
		1,
		domain.CustomDuration{
			Duration: time.Hour,
		},
		domain.CustomDuration{
			Duration: 10 * time.Minute,
		},
		domain.CustomDuration{
			Duration: 5 * time.Minute,
		},
		80,
	)

	expected := "[SUCCESS] 1 [01:00:00, 00:10:00, 00:05:00] HP:80\n"

	if line != expected {
		t.Fatalf("expected %q got %q", expected, line)
	}
}

func TestWriteReport(t *testing.T) {
	file, err := os.CreateTemp("", "report_test")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(file.Name())

	writer := &ReportWriter{
		file: file,
	}

	users := map[int]*domain.User{
		1: {
			ID:     1,
			Result: domain.ReportHeaderSuccess,
			Health: 100,
		},
	}

	writer.WriteReport(users)
}
