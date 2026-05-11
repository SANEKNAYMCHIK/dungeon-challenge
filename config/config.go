package config

import platformConfig "dungeon-challenge/platform/config"

type InputConfig struct {
	ConfigName string `yaml:"input/dungeon" env:"INPUT_DUNGEON" env-default:"config.json"`
	EventsName string `yaml:"input/events" env:"INPUT_EVENTS" env-default:"events.txt"`
}

type OutputConfig struct {
	OutputName string `yaml:"output/events" env:"OUTPUT_EVENTS" env-default:"output.txt"`
	ReportName string `yaml:"output/report" env:"OUTPUT_REPORT" env-default:"report.txt"`
}

type Config struct {
	Input  InputConfig  `yaml:"input"`
	Output OutputConfig `yaml:"output"`
}

func MustLoad(path string) Config {
	var cfg Config
	platformConfig.MustLoad(path, &cfg)
	return cfg
}
