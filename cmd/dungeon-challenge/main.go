package main

import (
	"dungeon-challenge/config"
	"dungeon-challenge/internal/controller/output"
	"dungeon-challenge/internal/controller/parser"
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
	defer ew.Close()

	dungeonParser := parser.NewDungeonParser(cfg.Input.ConfigName)
	dungeon, err := dungeonParser.ParseDungeon()
	if err != nil {
		log.Fatalf("error parsing dungeon: %v", err)
	}
	fmt.Println(dungeon.Duration)
}
