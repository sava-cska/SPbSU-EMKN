package accounts

import (
	"errors"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage"
	"github.com/sava-cska/SPbSU-EMKN/internal/utils"
	"github.com/sirupsen/logrus"
	"net/http"
)

func HandleAccountsLogin(logger *logrus.Logger, storage *storage.Storage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		logger.Debugf("HandleAccountsLogin - Called URI %s", request.RequestURI)

		var loginRequest loginRequest
		if err := utils.ParseBody(interface{}(&loginRequest), request); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte("Failed to parse request body"))
			return
		}

		isValid, err := ValidateUserCredentials(loginRequest.Login, loginRequest.Password, logger, storage)
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
func ValidateUserCredentials(login string, password string, logger *logrus.Logger, storage *storage.Storage) (bool, error) {
	origPassword, err := storage.UserDAO().GetPassword(login)
	if err != nil {
		logger.Errorf("Failed to read password for login %s: %s", login, err.Error())
		return false, errors.New("internal database error")
	}

	if password != origPassword {
		return false, nil
	} else {
		return true, nil
	}
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
