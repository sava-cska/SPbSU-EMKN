package accounts

import (
	"net/http"
	"strconv"
	"time"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/internal_data"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/notifier"
)

func HandleAccountsBeginChangePassword(request *BeginChangePasswordRequest,
	context *dependency.DependencyContext) (int, *BeginChangePasswordResponse) {
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
		if errEmail := context.Mailer.SendEmail([]string{request.Email}, notifier.BuildMessage(verificationCode,
			user.FirstName, user.LastName)); errEmail != nil {
			context.Logger.Errorf("BeginChangePassword: can't send email to %s, %s", request.Email, errEmail)
		}
	}()

	if errDB := context.Storage.ChangePasswordDao().Upsert(token, user.Login, time.Now().Add(internal_data.TokenTTL),
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
