package server

import "github.com/sava-cska/SPbSU-EMKN/internal/app/storage"

type Config struct {
	BindAddress string `toml:"bind_address"`
	LogLevel    string `toml:"log_level"`
	Storage     *storage.Config
}

// NewConfig Upsert new instance of config
func NewConfig() *Config {
	return &Config{
		BindAddress: ":8080",
		LogLevel:    "debug",
		Storage:     storage.NewConfig(),
	}
}
