package accounts

import (
	"encoding/json"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/notifier"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage"
	"github.com/sava-cska/SPbSU-EMKN/internal/utils"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

func HandleAccountsRevalidateRegistrationCredentials(logger *logrus.Logger, storage *storage.Storage, mailer *notifier.Mailer) http.HandlerFunc {
	resentCodeIn := 60 * time.Second
	tokenTTL := 30 * time.Minute
	verificationCodeLength := 6
	tokenLength := uint16(20)

	handleAccountsRevalidateRegistrationCredentials :=
		func(request *RevalidateRegistrationCredentialsRequest) (int, *RevalidateRegistrationCredentialsResponse) {
			user, _, _, err := storage.RegistrationDAO().FindRegistrationAndDelete(request.Token)
			if err != nil {
				return http.StatusInternalServerError, &RevalidateRegistrationCredentialsResponse{}
			}
			if storage.UserDAO().Exists(user.Login) {
				return http.StatusBadRequest, &RevalidateRegistrationCredentialsResponse{Errors: &ErrorsUnion{
					InvalidRegistrationRevalidation: &Error{}},
				}
			}
			token := generateToken(tokenLength)
			verificationCode := generateVerificationCode(verificationCodeLength)
			storage.RegistrationDAO().Upsert(
				token,
				&user,
				time.Now().Add(tokenTTL),
				verificationCode,
			)

			go func() {
				err := mailer.SendEmail([]string{user.Email}, buildMessage(verificationCode, user.FirstName, user.LastName))
				if err != nil {
					logger.Debugf(err.Error())
				}
			}()

			return http.StatusOK, &RevalidateRegistrationCredentialsResponse{
				Response: &RevalidateRegistrationCredentialsWrapper{
					RandomToken: token,
					ExpiresIn:   strconv.Itoa(int(resentCodeIn.Seconds())),
				},
			}
		}

	return func(writer http.ResponseWriter, request *http.Request) {
		logger.Debugf("HandleAccountsRevalidateRegistrationCredentials - Called URI %s", request.RequestURI)

		var revalidateRegistrationCredentialsRequest RevalidateRegistrationCredentialsRequest
		if errJSON := utils.ParseBody(interface{}(&revalidateRegistrationCredentialsRequest), request); errJSON != nil {
			utils.HandleError(logger,
				writer,
				http.StatusBadRequest,
				"Can't parse /accounts/revalidate_registration_credentials request.",
				errJSON)
			return
		}

		code, resp := handleAccountsRevalidateRegistrationCredentials(&revalidateRegistrationCredentialsRequest)

		respJSON, errRespJSON := json.Marshal(resp)
		if errRespJSON != nil {
			utils.HandleError(logger, writer, http.StatusInternalServerError, "Can't create JSON object from data.", errRespJSON)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(code)
		writer.Write(respJSON)
	}
}
