package models

import (
	"fmt"
	"html"
)

type Message struct {
	Subject   string
	Body      string
	Receivers []string
}

func BuildMessage(verificationCode string, firstName string, lastName string, receivers []string) Message {
	return Message{
		Subject: "Код подтверждения",
		Body: fmt.Sprintf("<html><body>Здравствуйте, %s %s!<br>Код подтверждения: <b>%s</b></body></html>",
			html.EscapeString(firstName),
			html.EscapeString(lastName),
			verificationCode),
		Receivers: receivers,
	}
}
