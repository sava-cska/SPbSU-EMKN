package accounts

import (
	"fmt"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"net/http"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage"
)

func HandleAccountsLogin(loginRequest *LoginRequest, context *dependency.DependencyContext) (int, *LoginResponse) {
	isValid, err := ValidateUserCredentials(loginRequest.Login, loginRequest.Password, context.Storage)
	if err != nil {
		context.Logger.Error("Failed to validate user credentials", err)
		return http.StatusInternalServerError, &LoginResponse{}
	}

	var statusCode int
	var response LoginResponse
	if !isValid {
		statusCode = http.StatusBadRequest
		response = LoginResponse{Errors: &ErrorsUnion{InvalidLoginOrPassword: &Error{}}}
	} else {
		statusCode = http.StatusOK
		response = LoginResponse{}
	}

	return statusCode, &response
}

// ValidateUserCredentials returns tuple (is credentials valid, error if internal error happened)
func ValidateUserCredentials(login string, password string, storage *storage.Storage) (bool, error) {
	origPassword, err := storage.UserDAO().GetPassword(login)
	if err != nil {
		return false, fmt.Errorf("failed to read password for login %s: %s", login, err.Error())
	}
	return password == origPassword, nil
}
