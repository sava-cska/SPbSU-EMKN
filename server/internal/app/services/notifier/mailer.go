package notifier

import (
	"encoding/base64"
	"fmt"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/models"
	"net/smtp"
	"strconv"
	"strings"
)

type Mailer interface {
	SendEmail(message models.Message) error
}

type mailerImpl struct {
	config *Config
	auth   smtp.Auth
	sender string
}

func New(config *Config, EmknCourseMail, EmknCoursePassword string) Mailer {
	return &mailerImpl{
		config: config,
		auth:   smtp.PlainAuth("", EmknCourseMail, EmknCoursePassword, config.MailerDaemon),
		sender: EmknCourseMail,
	}
}

func (mailer mailerImpl) SendEmail(message models.Message) error {
	msg := fmt.Sprintf(`From: %s
Subject: %s
Content-Type: text/html; charset="UTF-8"
Content-Transfer-Encoding: base64

%s`,
		mailer.sender,
		message.Subject,
		base64.StdEncoding.EncodeToString([]byte(message.Body)))
	return smtp.SendMail(mailer.getMailerDaemon(),
		mailer.auth,
		mailer.sender,
		message.Receivers,
		[]byte(msg))
}

func (mailer mailerImpl) getMailerDaemon() string {
	sb := strings.Builder{}
	sb.WriteString(mailer.config.MailerDaemon)
	sb.WriteString(":")
	sb.WriteString(strconv.Itoa(mailer.config.MailerDaemonPort))
	return sb.String()
}
