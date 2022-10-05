package accounts

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"net/mail"
	"strconv"
	"time"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/notifier"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage"
	"github.com/sava-cska/SPbSU-EMKN/internal/utils"

	"github.com/sirupsen/logrus"
)

func HandleAccountsRegister(logger *logrus.Logger, storage *storage.Storage, mailer *notifier.Mailer) http.HandlerFunc {
	validate := func(request *RegisterRequest) (int, *ErrorsUnion) {
		if len(request.Login) == 0 {
			return http.StatusBadRequest, &ErrorsUnion{
				IllegalLogin: &Error{},
			}
		}
		if len(request.Password) == 0 {
			return http.StatusBadRequest, &ErrorsUnion{
				IllegalPassword: &Error{},
			}
		}
		if _, err := mail.ParseAddress(request.Email); err != nil {
			return http.StatusBadRequest, &ErrorsUnion{
				IllegalEmail: &Error{},
			}
		}
		return http.StatusOK, nil
	}

	handleAccountsRegister := func(request *RegisterRequest) (int, *RegisterResponse) {
		if code, errors := validate(request); errors != nil {
			return code, &RegisterResponse{
				Errors: errors,
			}
		}
		if storage.UserDAO().Exists(request.Login) {
			return http.StatusBadRequest, &RegisterResponse{
				Errors: &ErrorsUnion{LoginIsNotAvailable: &Error{}},
			}
		}

		token := utils.GenerateToken()
		verificationCode := utils.GenerateVerificationCode()
		storage.RegistrationDAO().Upsert(
			token,
			request,
			time.Now().Add(utils.TokenTTL),
			verificationCode,
		)

		go func() {
			err := mailer.SendEmail([]string{request.Email}, buildMessage(verificationCode, request.FirstName, request.LastName))
			if err != nil {
				logger.Error("Can't send email to %s, %s", request.Email, err.Error())
			}
		}()

		return http.StatusOK, &RegisterResponse{
			Response: &RegisterWrapper{
				RandomToken: token,
				ExpiresIn:   strconv.Itoa(int(utils.ResentCodeIn.Seconds())),
			},
		}
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		logger.Debugf("HandleAccountsRegister - Called URI %s", request.RequestURI)

		var registerRequest RegisterRequest
		if errJSON := utils.ParseBody(interface{}(&registerRequest), request); errJSON != nil {
			utils.HandleError(logger, writer, http.StatusBadRequest, "Can't parse /accounts/register request.", errJSON)
			return
		}

		code, resp := handleAccountsRegister(&registerRequest)

		respJSON, errRespJSON := json.Marshal(resp)
		if errRespJSON != nil {
			utils.HandleError(logger, writer, http.StatusInternalServerError, "Can't create JSON object from data.", errRespJSON)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(code)
		writer.Write(respJSON)
	}
}

func buildMessage(verificationCode string, firstName string, lastName string) notifier.Message {
	return notifier.Message{
		Subject: "Код подтверждения",
		Body: fmt.Sprintf("<html><body>Здравствуйте, %s %s!<br>Код подтверждения: <b>%s</b></body></html>",
			html.EscapeString(firstName),
			html.EscapeString(lastName),
			verificationCode),
	}
}
