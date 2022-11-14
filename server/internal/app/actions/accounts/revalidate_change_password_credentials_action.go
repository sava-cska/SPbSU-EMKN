package accounts

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/internal_data"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/notifier"
)

// CreateOrder godoc
// @Summary Revalidate changing password
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param request body RevalidateChangePasswordCredentialsRequest true "Revalidate changing password"
// @Success 200 {object} RevalidateChangePasswordCredentialsResponse
// @Router /accounts/revalidate_change_password_credentials [post]
func HandleAccountsRevalidateChangePasswordCredentials(request *RevalidateChangePasswordCredentialsRequest,
	context *dependency.DependencyContext,
	_ ...any) (int, *RevalidateChangePasswordCredentialsResponse) {
	context.Logger.Debugf("RevalidateChangePasswordCredentials: start with token = %s", request.RandomToken)

	returnErr := func(statusCode int, reason string, err error) (int, *RevalidateChangePasswordCredentialsResponse) {
		context.Logger.Errorf("RevalidateChangePasswordCredentials: %s, %s", reason, err)
		return statusCode, &RevalidateChangePasswordCredentialsResponse{}
	}

	login, err := context.Storage.ChangePasswordDAO().FindTokenAndDelete(request.RandomToken)
	if err != nil {
		context.Logger.Errorf("RevalidateChangePasswordCredentials: can't find and delete record with token = %s",
			request.RandomToken)
		return http.StatusBadRequest, &RevalidateChangePasswordCredentialsResponse{
			Errors: &ErrorsUnion{
				InvalidChangePasswordRevalidation: &Error{},
			},
		}
	}

	newToken := internal_data.GenerateToken()
	verificationCode := notifier.GenerateVerificationCode()
	context.Logger.Debugf("RevalidateChangePasswordCredentials: token = %s, verificationCode = %s", newToken, verificationCode)

	user, errorUser := context.Storage.UserDAO().FindUserByLogin(login)
	if errorUser != nil {
		return returnErr(http.StatusInternalServerError, fmt.Sprintf("Can't find user with login = %s", login), errorUser)
	}
	context.Logger.Debugf("RevalidateChangePasswordCredentials: find user with login = %s and email = %s", user.Login, user.Email)

	go func() {
		if errEmail := context.Mailer.SendEmail([]string{user.Email}, notifier.BuildMessage(verificationCode,
			user.FirstName, user.LastName)); errEmail != nil {
			context.Logger.Errorf("RevalidateChangePasswordCredentials: can't send email to %s, %s", user.Email, errEmail)
		}
	}()

	if errDB := context.Storage.ChangePasswordDAO().Upsert(newToken, user.Login, time.Now().Add(internal_data.TokenTTL),
		verificationCode); errDB != nil {
		return returnErr(http.StatusInternalServerError, "Can't add record to change_password_base", errDB)
	}

	return http.StatusOK, &RevalidateChangePasswordCredentialsResponse{
		Response: &RevalidateChangePasswordCredentialsWrapper{
			RandomToken: newToken,
			ExpiresIn:   strconv.Itoa(int(internal_data.ResentCodeIn.Seconds())),
		},
	}
}
