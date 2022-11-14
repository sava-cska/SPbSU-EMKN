package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/core"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/config.toml", "path to config file")
}

// @title EMKN API
// @version 1.0
// @description This is a backend for EMKN app
// @contact.name API Support
// @contact.email https://t.me/Intellec2aI
// @host localhost:8080
// @BasePath /
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
