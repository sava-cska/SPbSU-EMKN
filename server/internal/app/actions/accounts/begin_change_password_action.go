package accounts

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/notifier"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage"
	"github.com/sava-cska/SPbSU-EMKN/internal/utils"
	"github.com/sirupsen/logrus"
)

func startChangePassword(logger *logrus.Logger, storage *storage.Storage, mailer *notifier.Mailer,
	email string) ([]byte, int, error) {
	user, errFindUser := storage.UserDAO().FindUser(email)
	if errFindUser != nil {
		response := ChangePwdResponse{Errors: &ErrorsUnion{IllegalEmail: &Error{}}}
		responseJSON, errJSON := json.Marshal(response)
		if errJSON != nil {
			logger.Error("Can't create JSON object from data.")
			return []byte{}, 0, errJSON
		}
		return responseJSON, http.StatusInternalServerError, nil
	}

	token := utils.GenerateToken()
	verificationCode := utils.GenerateVerificationCode()

	go func() {
		if errEmail := mailer.SendEmail([]string{email}, buildMessage(verificationCode,
			user.FirstName, user.LastName)); errEmail != nil {
			logger.Error("Can't send email to %s, %s", email, errEmail.Error())
		}
	}()

	if errDB := storage.ChangePasswordDao().Upsert(token, user.Login, time.Now().Add(utils.TokenTTL),
		verificationCode); errDB != nil {
		logger.Error("Can't add record to change_password_base.")
		return []byte{}, 0, errDB
	}

	successResponse := ChangePwdResponse{
		Response: &ChangePwdWrapper{
			Token:       token,
			TimeExpired: strconv.Itoa(int(utils.ResentCodeIn.Seconds()))}}
	successResponseJSON, errJSON := json.Marshal(successResponse)
	if errJSON != nil {
		logger.Error("Can't create JSON object from data.")
		return []byte{}, 0, errJSON
	}
	return successResponseJSON, http.StatusOK, nil
}

func HandleAccountsChangePwd(logger *logrus.Logger, storage *storage.Storage, mailer *notifier.Mailer) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		logger.Debugf("HandleAccountsChangePwd - Called URI %s", request.RequestURI)

		var changePwdRequest ChangePwdRequest
		if errJSON := utils.ParseBody(interface{}(&changePwdRequest), request); errJSON != nil {
			utils.HandleError(logger, writer, http.StatusBadRequest, "Can't parse /accounts/begin_change_password request.", errJSON)
			return
		}

		responseBody, returnCode, err := startChangePassword(logger, storage, mailer, changePwdRequest.Email)
		if err != nil {
			utils.HandleError(logger, writer, http.StatusInternalServerError, "Can't start process of changing password.", err)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(returnCode)
		writer.Write(responseBody)
	}
}
