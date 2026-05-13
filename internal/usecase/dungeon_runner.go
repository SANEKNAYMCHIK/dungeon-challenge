package usecase

import (
	"dungeon-challenge/internal/controller/output"
	"dungeon-challenge/internal/controller/parser"
	"dungeon-challenge/internal/domain"
	"errors"
	"fmt"
	"io"
	"log"
)

type DungeonRunner struct {
	dungeonInfo  domain.Dungeon
	eventsReader *parser.EventsParser
	outputWriter *output.Writer
	reportWriter *output.Writer
}

func NewDungeonRunner(dungeonInfo domain.Dungeon, eventsReader *parser.EventsParser, outputWriter *output.Writer, reportWriter *output.Writer) *DungeonRunner {
	return &DungeonRunner{
		dungeonInfo:  dungeonInfo,
		eventsReader: eventsReader,
		outputWriter: outputWriter,
		reportWriter: reportWriter,
	}
}

func (ds *DungeonRunner) Run() {
	for {
		event, err := ds.eventsReader.ReadEvent()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Printf("error reading event: %v", err)
			continue
		}
		fmt.Println(event)
	}
}
