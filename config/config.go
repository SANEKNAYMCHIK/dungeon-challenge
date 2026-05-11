package config

import "dungeon-challenge/platform/env"

type Config struct {
	ConfigName string `yaml:"input/dungeon" env:"INPUT_DUNGEON" env-default:"config.json"`
	EventsName string `yaml:"input/events" env:"INPUT_EVENTS" env-default:"events.txt"`
	OutputName string `yaml:"output/events" env:"OUTPUT_EVENTS" env-default:"output.txt"`
	ReportName string `yaml:"output/report" env:"OUTPUT_REPORT" env-default:"report.txt"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}
