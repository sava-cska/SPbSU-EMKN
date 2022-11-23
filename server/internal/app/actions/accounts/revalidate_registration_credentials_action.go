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
// @Summary Revalidate creds
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param request body RevalidateRegistrationCredentialsRequest true "Revalidate creds"
// @Success 200 {object} RevalidateRegistrationCredentialsResponse
// @Router /accounts/revalidate_registration_credentials [post]
func HandleAccountsRevalidateRegistrationCredentials(request *RevalidateRegistrationCredentialsRequest,
	context *dependency.DependencyContext,
	_ ...any) (int, *RevalidateRegistrationCredentialsResponse) {
	context.Logger.Debugf("RevalidateRegistrationCredentials: start with token = %s", request.Token)

	handleAccountsRevalidateRegistrationCredentials :=
		func(request *RevalidateRegistrationCredentialsRequest) (int, *RevalidateRegistrationCredentialsResponse) {
			user, _, _, err := context.Storage.RegistrationDAO().FindRegistrationAndDelete(request.Token)
			if err != nil {
				context.Logger.Errorf("RevalidateRegistrationCredentials: can't find and delete record in registration_base")
				return http.StatusBadRequest, &RevalidateRegistrationCredentialsResponse{Errors: &ErrorsUnion{
					InvalidRegistrationRevalidation: &Error{}},
				}
			}
			context.Logger.Debugf("RevalidateRegistrationCredentials: find user with login = %s", user.Login)
			if context.Storage.UserDAO().ExistsLogin(user.Login) {
				context.Logger.Errorf("RevalidateRegistrationCredentials: login = %s already exist", user.Login)
				return http.StatusBadRequest, &RevalidateRegistrationCredentialsResponse{Errors: &ErrorsUnion{
					InvalidRegistrationRevalidation: &Error{}},
				}
			}

			token := internal_data.GenerateToken()
			verificationCode := notifier.GenerateVerificationCode()
			context.Logger.Debugf("RevalidateRegistrationCredentials: token = %s, verificationCode = %s", token, verificationCode)

			context.Storage.RegistrationDAO().Upsert(
				token,
				&user,
				time.Now().Add(internal_data.TokenTTL),
				verificationCode,
			)

			go func() {
				err := context.EventQueue.AddMessage("Email", models.BuildMessage(verificationCode, user.FirstName,
					user.LastName, []string{user.Email}))
				if err != nil {
					context.Logger.Debugf("RevalidateRegistrationCredentials: can't send email to %s, %s", user.Email, err)
				}
			}()

			return http.StatusOK, &RevalidateRegistrationCredentialsResponse{
				Response: &RevalidateRegistrationCredentialsWrapper{
					RandomToken: token,
					ExpiresIn:   strconv.Itoa(int(internal_data.ResentCodeIn.Seconds())),
				},
			}
		}

	return handleAccountsRevalidateRegistrationCredentials(request)
}
