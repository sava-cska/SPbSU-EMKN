package accounts

import (
	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/notifier"
	"github.com/sava-cska/SPbSU-EMKN/internal/utils"
	"net/http"
)

func HandleRevalidateChangePasswordCredentials(request *RevalidateChangePasswordCredentialsRequest, context *dependency.DependencyContext) (int, *RevalidateChangePasswordCredentialsResponse) {
	returnErr := func(statusCode int, reason string, err error) (int, *RevalidateChangePasswordCredentialsResponse) {
		context.Logger.Error(reason, err)
		return statusCode, &RevalidateChangePasswordCredentialsResponse{}
	}

	verificationCode := utils.GenerateVerificationCode()
	login, tokenExists, err := context.Storage.ChangePasswordDao().UpdateVerificationCode(request.RandomToken, verificationCode)
	if err != nil {
		return returnErr(http.StatusInternalServerError, "Failed to update verification code", err)
	}

	var statusCode int
	var resp RevalidateChangePasswordCredentialsResponse
	if tokenExists {
		user, err := context.Storage.UserDAO().FindUserByLogin(login)
		if err != nil {
			return returnErr(http.StatusInternalServerError, "Failed to get user from database", err)
		}
		go func() {
			if errEmail := context.Mailer.SendEmail([]string{user.Email}, notifier.BuildMessage(verificationCode,
				user.FirstName, user.LastName)); errEmail != nil {
				context.Logger.Errorf("Can't send email to %s, %s", user.Email, errEmail.Error())
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

	return statusCode, &resp
}
