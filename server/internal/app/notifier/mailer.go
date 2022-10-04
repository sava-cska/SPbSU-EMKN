package notifier

import (
	"encoding/base64"
	"net/smtp"
	"strconv"
	"strings"
)

type Mailer struct {
	config *Config
	auth   smtp.Auth
	sender string
}

func New(config *Config, EmknCourseMail, EmknCoursePassword string) *Mailer {
	return &Mailer{
		config: config,
		auth:   smtp.PlainAuth("", EmknCourseMail, EmknCoursePassword, config.MailerDaemon),
		sender: EmknCourseMail,
	}
}

func (mailer Mailer) SendEmail(receivers []string, message Message) error {
	msg := "From: " + mailer.sender + "\r\n" +
		"Subject: " + message.Subject + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" +
		base64.StdEncoding.EncodeToString([]byte(message.Body))

	return smtp.SendMail(mailer.getMailerDaemon(),
		mailer.auth,
		mailer.sender,
		receivers,
		[]byte(msg))
}

func (mailer Mailer) getMailerDaemon() string {
	sb := strings.Builder{}
	sb.WriteString(mailer.config.MailerDaemon)
	sb.WriteString(":")
	sb.WriteString(strconv.Itoa(mailer.config.MailerDaemonPort))
	return sb.String()
}
