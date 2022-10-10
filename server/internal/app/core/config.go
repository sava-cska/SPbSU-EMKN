package core

import (
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/notifier"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage"
)

type Config struct {
	BindAddress    string `toml:"bind_address"`
	LogLevel       string `toml:"log_level"`
	StorageConfig  *storage.Config
	NotifierConfig *notifier.Config
}

func NewConfig() *Config {
	return &Config{
		BindAddress:    ":8080",
		LogLevel:       "debug",
		StorageConfig:  storage.NewConfig(),
		NotifierConfig: notifier.NewConfig(),
	}
}
