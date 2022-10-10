package accounts

import (
	"encoding/json"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/notifier"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage"
	"github.com/sava-cska/SPbSU-EMKN/internal/utils"
	"github.com/sirupsen/logrus"
	"net/http"
)

func HandleRevalidateChangePasswordCredentials(logger *logrus.Logger, storage *storage.Storage, mailer *notifier.Mailer) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		logger.Debugf("HandleRevalidateChangePasswordCredentials - Called URI %s", request.RequestURI)

		var parsedRequest = RevalidateChangePasswordCredentialsRequest{}
		if err := utils.ParseBody(interface{}(&parsedRequest), request); err != nil {
			utils.HandleError(logger, writer, http.StatusBadRequest, "Failed to parse revalidate_change_password request body", err)
			return
		}

		verificationCode := utils.GenerateVerificationCode()
		login, tokenExists, err := storage.ChangePasswordDao().UpdateVerificationCode(parsedRequest.RandomToken, verificationCode)
		if err != nil {
			utils.HandleError(logger, writer, http.StatusInternalServerError, "Failed to update verification code", err)
			return
		}

		var statusCode int
		var resp RevalidateChangePasswordCredentialsResponse
		if tokenExists {
			user, err := storage.UserDAO().FindUserByLogin(login)
			if err != nil {
				utils.HandleError(logger, writer, http.StatusInternalServerError, "Failed to get user from database", err)
				return
			}
			go func() {
				if errEmail := mailer.SendEmail([]string{user.Email}, notifier.BuildMessage(verificationCode,
					user.FirstName, user.LastName)); errEmail != nil {
					logger.Errorf("Can't send email to %s, %s", user.Email, errEmail.Error())
				}
			}()

			resp = RevalidateChangePasswordCredentialsResponse{}
			statusCode = http.StatusOK
		} else {
			resp = RevalidateChangePasswordCredentialsResponse{
				Errors: &ErrorsUnion{
					InvalidChangePasswordRevalidation: &Error{},
				},
			}
			statusCode = http.StatusBadRequest
		}

		body, err := json.Marshal(resp)
		if err != nil {
			utils.HandleError(logger, writer, http.StatusInternalServerError, "Failed to write revalidate_change_password response", err)
			return
		}

		writer.WriteHeader(statusCode)
		_, _ = writer.Write(body)
		writer.Header().Set("Content-Type", "application/json")
	}
}
