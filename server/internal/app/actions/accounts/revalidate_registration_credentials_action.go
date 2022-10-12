package accounts

import (
	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/notifier"
	"net/http"
	"strconv"
	"time"

	"github.com/sava-cska/SPbSU-EMKN/internal/utils"
)

func HandleAccountsRevalidateRegistrationCredentials(request *RevalidateRegistrationCredentialsRequest,
	context *dependency.DependencyContext) (int, *RevalidateRegistrationCredentialsResponse) {
	handleAccountsRevalidateRegistrationCredentials :=
		func(request *RevalidateRegistrationCredentialsRequest) (int, *RevalidateRegistrationCredentialsResponse) {
			user, _, _, err := context.Storage.RegistrationDAO().FindRegistrationAndDelete(request.Token)
			if err != nil {
				return http.StatusInternalServerError, &RevalidateRegistrationCredentialsResponse{}
			}
			if context.Storage.UserDAO().ExistsLogin(user.Login) {
				return http.StatusBadRequest, &RevalidateRegistrationCredentialsResponse{Errors: &ErrorsUnion{
					InvalidRegistrationRevalidation: &Error{}},
				}
			}
			token := utils.GenerateToken()
			verificationCode := utils.GenerateVerificationCode()
			context.Storage.RegistrationDAO().Upsert(
				token,
				&user,
				time.Now().Add(utils.TokenTTL),
				verificationCode,
			)

			go func() {
				err := context.Mailer.SendEmail([]string{user.Email}, notifier.BuildMessage(verificationCode, user.FirstName, user.LastName))
				if err != nil {
					context.Logger.Debugf(err.Error())
				}
			}()

			return http.StatusOK, &RevalidateRegistrationCredentialsResponse{
				Response: &RevalidateRegistrationCredentialsWrapper{
					RandomToken: token,
					ExpiresIn:   strconv.Itoa(int(utils.ResentCodeIn.Seconds())),
				},
			}
		}

	return handleAccountsRevalidateRegistrationCredentials(request)
}
