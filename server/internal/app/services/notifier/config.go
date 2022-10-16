package notifier

type Config struct {
	MailerDaemon     string `toml:"mailer_daemon"`
	MailerDaemonPort int    `toml:"mailer_daemon_port"`
}

func NewConfig() *Config {
	return &Config{}
}
