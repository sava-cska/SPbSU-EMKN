package accounts

import (
	"fmt"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage"
	"net/http"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
)

func HandleAccountsLogin(loginRequest *LoginRequest, context *dependency.DependencyContext, _ ...any) (int, *LoginResponse) {
	context.Logger.Debugf("Login: start with login = %s, password = %s", loginRequest.Login, loginRequest.Password)

	user, err := context.Storage.UserDAO().FindUserByLogin(loginRequest.Login)
	if err != nil {
		return http.StatusInternalServerError, &LoginResponse{}
	}
	isValid, err := ValidateUserCredentials(loginRequest.Login, loginRequest.Password, context.Storage)
	if err != nil {
		context.Logger.Errorf("Login: failed to validate user credentials, %s", err)
		return http.StatusInternalServerError, &LoginResponse{}
	}

	if !isValid {
		context.Logger.Errorf("Login: incorrect credentials")
		return http.StatusBadRequest, &LoginResponse{
			Errors: &ErrorsUnion{
				InvalidLoginOrPassword: &Error{},
			},
		}
	}
	return http.StatusOK, &LoginResponse{
		ProfileId: &user.ProfileId,
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
