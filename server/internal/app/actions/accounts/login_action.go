package accounts

import (
	"fmt"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage"
	"github.com/sava-cska/SPbSU-EMKN/internal/utils"
	"github.com/sirupsen/logrus"
	"net/http"
)

func HandleAccountsLogin(logger *logrus.Logger, storage *storage.Storage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		logger.Debugf("HandleAccountsLogin - Called URI %s", request.RequestURI)

		var loginRequest LoginRequest
		if err := utils.ParseBody(interface{}(&loginRequest), request); err != nil {
			utils.HandleError(logger, writer, http.StatusBadRequest, "Failed to parse request body", err)
			return
		}

		isValid, err := ValidateUserCredentials(loginRequest.Login, loginRequest.Password, storage)
		if err != nil {
			utils.HandleError(logger, writer, http.StatusInternalServerError, "Failed to validate user credentials", err)
			return
		}

		if !isValid {
			writer.WriteHeader(http.StatusUnauthorized)
		} else {
			writer.WriteHeader(http.StatusOK)
		}
	}
}

// ValidateUserCredentials returns tuple (is credentials valid, error if internal error happened)
func ValidateUserCredentials(login string, password string, storage *storage.Storage) (bool, error) {
	origPassword, err := storage.UserDAO().GetPassword(login)
	if err != nil {
		return false, fmt.Errorf("failed to read password for login %s: %s", login, err.Error())
	}
	return password == origPassword, nil
}
