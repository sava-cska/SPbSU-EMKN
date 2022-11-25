package accounts

import (
	"github.com/sava-cska/SPbSU-EMKN/internal/app/models"
	"net/http"
	"strconv"
	"time"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/internal_data"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/notifier"
)

// CreateOrder godoc
// @Summary Begin change password
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param request body BeginChangePasswordRequest true "Begin change password"
// @Success 200 {object} BeginChangePasswordResponse
// @Router /accounts/begin_change_password [post]
func HandleAccountsBeginChangePassword(request *BeginChangePasswordRequest,
	context *dependency.DependencyContext,
	_ ...any) (int, *BeginChangePasswordResponse) {
	context.Logger.Debugf("BeginChangePassword: start with email = %s", request.Email)

	user, errFindUser := context.Storage.UserDAO().FindUser(request.Email)
	if errFindUser != nil {
		context.Logger.Errorf("BeginChangePassword: can't find user in user_base, %s", errFindUser)
		return http.StatusBadRequest, &BeginChangePasswordResponse{
			Errors: &ErrorsUnion{
				IllegalEmail: &Error{},
			},
		}
	}

	token := internal_data.GenerateToken()
	verificationCode := notifier.GenerateVerificationCode()
	context.Logger.Debugf("BeginChangePassword: token = %s, verificationCode = %s", token, verificationCode)

	go func() {
		if err := context.EventQueue.AddMessage("Email", models.BuildMessage(verificationCode,
			user.FirstName, user.LastName, []string{request.Email})); err != nil {
			context.Logger.Errorf("BeginChangePassword: can't send email to %s", request.Email)
		}
	}()

	if errDB := context.Storage.ChangePasswordDAO().UpsertChangePasswordData(token, user.Login, time.Now().Add(internal_data.TokenTTL),
		verificationCode); errDB != nil {
		context.Logger.Errorf("BeginChangePassword: can't add record to change_password_base, %s", errDB)
		return http.StatusInternalServerError, &BeginChangePasswordResponse{}
	}

	return http.StatusOK, &BeginChangePasswordResponse{
		Response: &BeginChangePasswordWrapper{
			Token:       token,
			TimeExpired: strconv.Itoa(int(internal_data.ResentCodeIn.Seconds())),
		},
	}
}
