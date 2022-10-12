package accounts

import (
	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"net/http"
	"time"

	"github.com/sava-cska/SPbSU-EMKN/internal/utils"
)

func HandleValidateChangePassword(request *ValidateChangePasswordRequest, context *dependency.DependencyContext) (int, *ValidateChangePasswordResponse) {

	returnErr := func(statusCode int, reason string, err error) (int, *ValidateChangePasswordResponse) {
		context.Logger.Error(reason, err)
		return statusCode, &ValidateChangePasswordResponse{}
	}

	correctVerificationCode, expiresAt, err := context.Storage.ChangePasswordDao().GetVerificationCodeInfo(request.RandomToken)
	if err != nil {
		return returnErr(http.StatusInternalServerError, "Failed to get verification code", err)
	}
	if correctVerificationCode == "" {
		return returnErr(http.StatusBadRequest, "Failed to find random token", nil)
	}

	errors := validateVerificationCode(correctVerificationCode, request.VerificationCode, expiresAt)

	var responseBody ValidateChangePasswordResponse
	var statusCode int
	if errors == nil {
		token := utils.GenerateToken()
		err = context.Storage.ChangePasswordDao().SetChangePasswordToken(request.RandomToken,
			time.Now().Add(utils.TokenTTL), token)
		if err != nil {
			return returnErr(http.StatusInternalServerError, "Failed to store changePasswordToken", err)
		}
		responseBody = ValidateChangePasswordResponse{
			ChangePasswordToken: token,
		}
		statusCode = http.StatusOK
	} else {
		responseBody = ValidateChangePasswordResponse{
			Errors: errors,
		}
		statusCode = http.StatusBadRequest
	}

	return statusCode, &responseBody
}

func validateVerificationCode(correctVerificationCode string, verificationCode string, expiresAt *time.Time) *ErrorsUnion {
	if expiresAt.Before(time.Now()) {
		return &ErrorsUnion{ChangePasswordExpired: &Error{}}
	}

	if correctVerificationCode != verificationCode {
		return &ErrorsUnion{InvalidCode: &Error{}}
	}

	return nil
}
