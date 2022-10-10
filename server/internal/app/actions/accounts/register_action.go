package accounts

import (
	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/models"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/notifier"
	"net/http"
	"net/mail"
	"strconv"
	"time"

	"github.com/sava-cska/SPbSU-EMKN/internal/utils"
)

func HandleAccountsRegister(request *RegisterRequest, context *dependency.DependencyContext) (int, *RegisterResponse) {
	validate := func(request *RegisterRequest) (int, *ErrorsUnion) {
		if len(request.Login) == 0 {
			return http.StatusBadRequest, &ErrorsUnion{
				IllegalLogin: &Error{},
			}
		}
		if len(request.Password) == 0 {
			return http.StatusBadRequest, &ErrorsUnion{
				IllegalPassword: &Error{},
			}
		}
		if _, err := mail.ParseAddress(request.Email); err != nil {
			return http.StatusBadRequest, &ErrorsUnion{
				IllegalEmail: &Error{},
			}
		}
		return http.StatusOK, nil
	}

	handleAccountsRegister := func(request *RegisterRequest) (int, *RegisterResponse) {
		if code, errors := validate(request); errors != nil {
			return code, &RegisterResponse{
				Errors: errors,
			}
		}
		if context.Storage.UserDAO().ExistsLogin(request.Login) {
			return http.StatusBadRequest, &RegisterResponse{
				Errors: &ErrorsUnion{LoginIsNotAvailable: &Error{}},
			}
		}
		if context.Storage.UserDAO().ExistsEmail(request.Email) {
			return http.StatusBadRequest, &RegisterResponse{
				Errors: &ErrorsUnion{EmailIsNotAvailable: &Error{}},
			}
		}

		token := utils.GenerateToken()
		verificationCode := utils.GenerateVerificationCode()
		context.Storage.RegistrationDAO().Upsert(
			token,
			&models.User{
				Login:     request.Login,
				Password:  request.Password,
				Email:     request.Email,
				FirstName: request.FirstName,
				LastName:  request.LastName,
			},
			time.Now().Add(utils.TokenTTL),
			verificationCode,
		)

		go func() {
			err := context.Mailer.SendEmail([]string{request.Email}, notifier.BuildMessage(verificationCode, request.FirstName, request.LastName))
			if err != nil {
				context.Logger.Error("Can't send email to %s, %s", request.Email, err.Error())
			}
		}()

		return http.StatusOK, &RegisterResponse{
			Response: &RegisterWrapper{
				RandomToken: token,
				ExpiresIn:   strconv.Itoa(int(utils.ResentCodeIn.Seconds())),
			},
		}
	}

	return handleAccountsRegister(request)
}
