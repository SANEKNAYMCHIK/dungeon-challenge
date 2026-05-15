package domain

import (
	"testing"
	"time"
)

func TestAverageDuration(t *testing.T) {
	durations := []CustomDuration{
		{Duration: 10 * time.Second},
		{Duration: 20 * time.Second},
	}

	avg := AverageDuration(durations)

	if avg.Duration != 15*time.Second {
		t.Fatalf("expected 15s, got %v", avg.Duration)
	}
}

func TestDurationString(t *testing.T) {
	d := CustomDuration{
		Duration: time.Hour + 2*time.Minute + 3*time.Second,
	}

	expected := "01:02:03"

	if d.String() != expected {
		t.Fatalf("expected %s, got %s", expected, d.String())
	}
}
