package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"os"
)

func ParseBody(value any, request *http.Request) error {
	data, errReq := io.ReadAll(request.Body)
	if errReq != nil {
		return errReq
	}
	if errJSON := json.Unmarshal(data, value); errJSON != nil {
		return errJSON
	}

	return nil
}

func SendEmail(email, verificationCode, firstName, lastName string) error {
	from := os.Getenv("EMKN_COURSE_MAIL")
	password := os.Getenv("EMKN_COURSE_PASSWORD")

	msg := "To: " + email + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + "Код подтверждения" + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n" +
		"\r\n" +
		base64.StdEncoding.EncodeToString([]byte(
			fmt.Sprintf("<html><body>Здравствуйте, \"%s %s\"!<br>Код подтверждения: <b>%s</b></body></html>",
				firstName,
				lastName,
				verificationCode)))

	auth := smtp.PlainAuth("", from, password, "smtp.yandex.ru")
	return smtp.SendMail("smtp.yandex.ru:25", auth, from, []string{email}, []byte(msg))
}
