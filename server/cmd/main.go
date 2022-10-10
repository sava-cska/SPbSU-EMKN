package main

import (
	"flag"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/core"
	"log"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/config.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := core.NewConfig()
	_, err := toml.DecodeFile(configPath, config)

	if err != nil {
		log.Fatal(err)
	}

	s := core.New(config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
