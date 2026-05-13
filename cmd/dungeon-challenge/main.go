package main

import (
	"dungeon-challenge/config"
	"dungeon-challenge/internal/controller/output"
	"dungeon-challenge/internal/controller/parser"
	"dungeon-challenge/internal/usecase"
	"flag"
	"fmt"
	"log"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config/config.yaml", "server configuration file")
	fmt.Println(configPath)
	cfg := config.MustLoad(configPath)
	fmt.Println(cfg)
	ew := output.MustMakeWriter(cfg.Output.OutputName)
	defer func() {
		if err := ew.Close(); err != nil {
			log.Printf("failed to close output file: %v", err)
		}
	}()

	dungeonParser := parser.NewDungeonParser(cfg.Input.ConfigName)
	dungeon, err := dungeonParser.ParseDungeon()
	if err != nil {
		log.Fatalf("error parsing dungeon: %v", err)
	}
	eventsParser, err := parser.NewEventsParser(cfg.Input.EventsName)
	if err != nil {
		log.Fatalf("error creating events parser: %v", err)
	}
	defer func() {
		if err := eventsParser.Close(); err != nil {
			log.Printf("failed to close events file: %v", err)
		}
	}()

	dungeonRunner := usecase.NewDungeonRunner(dungeon, eventsParser, ew, ew)
	dungeonRunner.Run()
}
