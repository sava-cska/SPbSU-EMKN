package notifier

import (
	"fmt"
	"html"
)

type Message struct {
	Subject string
	Body    string
}

func BuildMessage(verificationCode string, firstName string, lastName string) Message {
	return Message{
		Subject: "Код подтверждения",
		Body: fmt.Sprintf("<html><body>Здравствуйте, %s %s!<br>Код подтверждения: <b>%s</b></body></html>",
			html.EscapeString(firstName),
			html.EscapeString(lastName),
			verificationCode),
	}
}
