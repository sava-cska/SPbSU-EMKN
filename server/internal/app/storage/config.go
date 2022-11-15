package storage

type Config struct {
	DatabaseURL string `toml:"database_url"`
}

func (c *Config) Auth(user string, password string) {
	c.DatabaseURL += " user=" + user + " password=" + password
}

func NewConfig() *Config {
	return &Config{}
}
