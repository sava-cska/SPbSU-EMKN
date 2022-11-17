package accounts

import (
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/pwd_hasher"
	"net/http"
	"time"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
)

// CreateOrder godoc
// @Summary Validate email code
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param request body ValidateEmailRequest true "Validate email code"
// @Success 200 {object} ValidateEmailResponse
// @Router /accounts/validate_email [post]
func HandleAccountsValidateEmail(request *ValidateEmailRequest, context *dependency.DependencyContext, _ ...any) (int, *ValidateEmailResponse) {
	context.Logger.Debugf("ValidateEmail: start with verificationCode = %s", request.VerificationCode)

	currentTime := time.Now()

	user, expireTime, verificationCodeDB, err := context.Storage.RegistrationDAO().FindRegistration(request.Token)
	if err != nil {
		context.Logger.Errorf("ValidateEmail: can't find record with token = %s in registration_base",
			request.Token)
		return http.StatusBadRequest, &ValidateEmailResponse{}
	}

	if expireTime.Before(currentTime) {
		context.Logger.Errorf("ValidateEmail: time expired")
		return http.StatusBadRequest, &ValidateEmailResponse{Errors: &ErrorsUnion{RegistrationExpired: &Error{}}}
	}
	if request.VerificationCode != verificationCodeDB {
		context.Logger.Errorf("ValidateEmail: verification code isn't correct, correct verification code = %s", verificationCodeDB)
		return http.StatusBadRequest, &ValidateEmailResponse{Errors: &ErrorsUnion{InvalidCode: &Error{}}}
	}

	user.Password, err = pwd_hasher.HashPassword(user.Password)
	if err != nil {
		context.Logger.Errorf("Failed to hash password: %s", err.Error())
		return http.StatusInternalServerError, &ValidateEmailResponse{}
	}

	if errAdd := context.Storage.UserDAO().AddUser(&user); errAdd != nil {
		context.Logger.Errorf("ValidateEmail: can't add user in user_base")
		return http.StatusInternalServerError, &ValidateEmailResponse{}
	}

	return http.StatusOK, &ValidateEmailResponse{}
}
