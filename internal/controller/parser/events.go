package parser

import (
	"bufio"
	"dungeon-challenge/internal/domain"
	"dungeon-challenge/internal/dto"
	"io"
	"os"
)

type EventsParser struct {
	EventsFile *bufio.Scanner
	file       *os.File
}

func NewEventsParser(filename string) (*EventsParser, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	return &EventsParser{
		EventsFile: scanner,
		file:       file,
	}, nil
}

func (ep *EventsParser) Close() error {
	return ep.file.Close()
}

func (ep *EventsParser) ReadEvent() (domain.Event, error) {
	if ep.EventsFile.Scan() {
		data := ep.EventsFile.Text()
		ReadyEvent, err := dto.ToEventDomain(data)
		if err != nil {
			return domain.Event{}, err
		}
		return ReadyEvent, nil
	}
	return domain.Event{}, io.EOF
}
