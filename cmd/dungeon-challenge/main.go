package main

import (
	"dungeon-challenge/config"
	"flag"
	"fmt"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config/config.yaml", "server configuration file")
	fmt.Println(configPath)
	cfg := config.MustLoad(configPath)
	fmt.Println(cfg)
}
