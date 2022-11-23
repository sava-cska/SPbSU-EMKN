package core

import (
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/event_queue"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/notifier"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage"
)

type Config struct {
	BindAddress      string `toml:"bind_address"`
	LogLevel         string `toml:"log_level"`
	StorageConfig    *storage.Config
	NotifierConfig   *notifier.Config
	EventQueueConfig *event_queue.Config
}

func NewConfig() *Config {
	return &Config{
		BindAddress:      ":8080",
		LogLevel:         "debug",
		StorageConfig:    storage.NewConfig(),
		NotifierConfig:   notifier.NewConfig(),
		EventQueueConfig: event_queue.NewConfig(),
	}
}
