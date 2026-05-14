package output

import (
	"dungeon-challenge/internal/domain"
	"fmt"
	"log"
	"os"
)

type EventWriter struct {
	file *os.File
}

func getOutputLine(eventID domain.EventType, lineTime domain.CustomTime, userID int, params string) string {
	if params == "" {
		return fmt.Sprintf(templates[eventID], lineTime, userID)
	} else {
		return fmt.Sprintf(templates[eventID], lineTime, userID, params)
	}
}

func (ew *EventWriter) WriteEvent(eventID domain.EventType, event domain.Event) (int, error) {
	lineTime, userID, params := event.Time, event.User, event.Param
	outputLine := getOutputLine(eventID, lineTime, userID, params)
	n, err := ew.file.WriteString(outputLine)
	return n, err
}

func (ew *EventWriter) WriteImpossibleMove(eventID domain.EventType, event domain.Event, params string) (int, error) {
	lineTime, userID := event.Time, event.User
	outputLine := getOutputLine(eventID, lineTime, userID, params)
	n, err := ew.file.WriteString(outputLine)
	return n, err
}

func (ew *EventWriter) WriteDeadUser(eventID domain.EventType, event domain.Event) (int, error) {
	lineTime, userID := event.Time, event.User
	outputLine := getOutputLine(eventID, lineTime, userID, "")
	n, err := ew.file.WriteString(outputLine)
	return n, err
}

func (ew *EventWriter) Close() error {
	return ew.file.Close()
}

func MustMakeWriter(filename string) *EventWriter {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	return &EventWriter{file: file}
}
