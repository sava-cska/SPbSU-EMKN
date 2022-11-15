package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/core"
)

var (
	configPath   string
	psqlPassword string
	psqlUser     string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/config.toml", "path to config file")
	flag.StringVar(&psqlPassword, "psql-password", "", "password for psql auth")
	flag.StringVar(&psqlUser, "psql-user", "", "user for psql auth")
}

// @title EMKN API
// @version 1.0
// @description This is a backend for EMKN app
// @contact.name API Support
// @contact.email https://t.me/Intellec2aI
// @host 51.250.98.212:8080
// @BasePath /
func main() {
	flag.Parse()

	config := core.NewConfig()
	_, err := toml.DecodeFile(configPath, config)

	if err != nil {
		log.Fatal(err)
	}
	config.StorageConfig.Auth(psqlUser, psqlPassword)

	s := core.New(config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
