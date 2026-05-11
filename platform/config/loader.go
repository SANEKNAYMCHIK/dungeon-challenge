package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

func MustLoad(path string, cfg any) {
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		log.Fatalf("cannot read config: %v", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		log.Fatalf("cannot read env: %v", err)
	}
}
