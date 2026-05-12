package main

import (
	"dungeon-challenge/config"
	"dungeon-challenge/platform/output"
	"flag"
	"fmt"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config/config.yaml", "server configuration file")
	fmt.Println(configPath)
	cfg := config.MustLoad(configPath)
	fmt.Println(cfg)
	ew := output.MustMakeWriter(cfg.Output.OutputName)
	defer ew.Close()
	_ = ew

}
