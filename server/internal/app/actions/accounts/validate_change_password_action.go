package accounts

import (
	"net/http"
	"time"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/internal_data"
)

func HandleAccountsValidateChangePassword(request *ValidateChangePasswordRequest,
	context *dependency.DependencyContext) (int, *ValidateChangePasswordResponse) {
	context.Logger.Debugf("ValidateChangePassword: start with verificationCode = %s", request.VerificationCode)

	currentTime := time.Now()

	returnErr := func(statusCode int, reason string, err error) (int, *ValidateChangePasswordResponse) {
		context.Logger.Errorf("ValidateChangePassword: %s, %s", reason, err)
		return statusCode, &ValidateChangePasswordResponse{}
	}

	correctVerificationCode, expiresAt, err := context.Storage.ChangePasswordDAO().GetVerificationCodeInfo(request.RandomToken)
	if err != nil {
		return returnErr(http.StatusInternalServerError, "Failed to get verification code from change_password_base", err)
	}
	if correctVerificationCode == "" {
		return returnErr(http.StatusBadRequest, "Failed to find random token in change_password_base", nil)
	}

	if expiresAt.Before(currentTime) {
		context.Logger.Errorf("ValidateChangePassword: time expired")
		return http.StatusBadRequest, &ValidateChangePasswordResponse{
			Errors: &ErrorsUnion{
				ChangePasswordExpired: &Error{},
			},
		}
	}

	if correctVerificationCode != request.VerificationCode {
		context.Logger.Errorf("ValidateChangePassword: verification code isn't correct, correct verification code = %s",
			correctVerificationCode)
		return http.StatusBadRequest, &ValidateChangePasswordResponse{
			Errors: &ErrorsUnion{
				InvalidCode: &Error{},
			},
		}
	}

	token := internal_data.GenerateToken()
	context.Logger.Debugf("ValidateChangePassword: token = %s", token)

	err = context.Storage.ChangePasswordDAO().SetChangePasswordToken(request.RandomToken,
		time.Now().Add(internal_data.TokenTTL), token)
	if err != nil {
		return returnErr(http.StatusInternalServerError, "Failed to store changePasswordToken in change_password_base", err)
	}
	return http.StatusOK, &ValidateChangePasswordResponse{ChangePasswordToken: token}
}
