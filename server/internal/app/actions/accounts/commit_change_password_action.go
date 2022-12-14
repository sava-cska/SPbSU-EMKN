package accounts

import (
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/pwd_hasher"
	"net/http"
	"time"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/internal_data"
)

// CreateOrder godoc
// @Summary Commit change password
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param request body CommitChangePasswordRequest true "Commit change password"
// @Success 200 {object} CommitChangePasswordResponse
// @Router /accounts/commit_change_password [post]
func HandleAccountsCommitChangePassword(request *CommitChangePasswordRequest,
	context *dependency.DependencyContext,
	_ ...any) (int, *CommitChangePasswordResponse) {
	context.Logger.Debugf("CommitChangePassword: start with new_password = %s", request.NewPassword)

	currentTime := time.Now()

	login, changePasswordExpiredTime, errDB := context.Storage.ChangePasswordDAO().FindPwdToken(request.ChangePasswordToken)
	if errDB != nil {
		context.Logger.Errorf("CommitChangePassword: can't find record with token = %s in change_password_base, %s",
			request.ChangePasswordToken, errDB)
		return http.StatusBadRequest, &CommitChangePasswordResponse{}
	}

	if changePasswordExpiredTime.Before(currentTime) {
		context.Logger.Errorf("CommitChangePassword: time for login = %s expired", login)
		return http.StatusBadRequest, &CommitChangePasswordResponse{
			Errors: &ErrorsUnion{
				ChangePasswordExpired: &Error{},
			},
		}
	}

	if !internal_data.ValidatePassword(request.NewPassword) {
		context.Logger.Errorf("CommitChangePassword: invalid password = %s", request.NewPassword)
		return http.StatusBadRequest, &CommitChangePasswordResponse{
			Errors: &ErrorsUnion{
				IllegalPassword: &Error{},
			},
		}
	}

	hashedPassword, err := pwd_hasher.HashPassword(request.NewPassword)
	if err != nil {
		context.Logger.Errorf("Failed to hash password: %s", err.Error())
		return http.StatusInternalServerError, &CommitChangePasswordResponse{}
	}

	if errUpdate := context.Storage.UserDAO().UpdatePassword(login, hashedPassword); errUpdate != nil {
		context.Logger.Errorf("CommitChangePassword: can't update password for login = %s in user_base", login, errUpdate)
		return http.StatusInternalServerError, &CommitChangePasswordResponse{}
	}

	return http.StatusOK, &CommitChangePasswordResponse{}
}
