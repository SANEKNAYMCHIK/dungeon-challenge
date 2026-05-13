package output

import (
	"dungeon-challenge/internal/domain"
	"fmt"
	"log"
	"os"
)

// TODO: EventWrite struct instead of Writer to use custom formatting of output information

type Writer struct {
	file *os.File
}

func getOutputLine(eventID domain.EventType, lineTime domain.CustomTime, userID int, params string) string {
	if params == "" {
		return fmt.Sprintf(templates[eventID], lineTime, userID)
	} else {
		return fmt.Sprintf(templates[eventID], lineTime, userID, params)
	}
}

func (ew *Writer) Write(eventID domain.EventType, lineTime domain.CustomTime, userID int, params string) (int, error) {
	outputLine := getOutputLine(eventID, lineTime, userID, params)
	n, err := ew.file.WriteString(outputLine)
	return n, err
}

func (ew *Writer) Close() error {
	return ew.file.Close()
}

func MustMakeWriter(filename string) *Writer {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	return &Writer{file: file}
}
