package event_queue

type Config struct {
	Address string `toml:"event_queue_address"`
	Port    string `toml:"event_queue_port"`
}

func NewConfig() *Config {
	return &Config{}
}
