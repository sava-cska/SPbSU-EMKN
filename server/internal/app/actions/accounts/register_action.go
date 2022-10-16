package accounts

import (
	"net/http"
	"strconv"
	"time"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/models"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/internal_data"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/services/notifier"
)

func HandleAccountsRegister(request *RegisterRequest, context *dependency.DependencyContext) (int, *RegisterResponse) {
	context.Logger.Debugf("Register: %s %s with email = %s, login = %s, password = %s", request.FirstName, request.LastName,
		request.Email, request.Login, request.Password)

	validate := func(request *RegisterRequest) (int, *ErrorsUnion) {
		if !internal_data.ValidateLogin(request.Login) {
			context.Logger.Errorf("Register: invalid login = %s", request.Login)
			return http.StatusBadRequest, &ErrorsUnion{
				IllegalLogin: &Error{},
			}
		}
		if !internal_data.ValidatePassword(request.Password) {
			context.Logger.Errorf("Register: invalid password = %s", request.Password)
			return http.StatusBadRequest, &ErrorsUnion{
				IllegalPassword: &Error{},
			}
		}
		if !internal_data.ValidateEmail(request.Email) {
			context.Logger.Errorf("Register: invalid email = %s", request.Email)
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
			context.Logger.Errorf("Register: login = %s already exist", request.Login)
			return http.StatusBadRequest, &RegisterResponse{
				Errors: &ErrorsUnion{LoginIsNotAvailable: &Error{}},
			}
		}
		if context.Storage.UserDAO().ExistsEmail(request.Email) {
			context.Logger.Errorf("Register: email = %s already exist", request.Email)
			return http.StatusBadRequest, &RegisterResponse{
				Errors: &ErrorsUnion{EmailIsNotAvailable: &Error{}},
			}
		}

		token := internal_data.GenerateToken()
		verificationCode := notifier.GenerateVerificationCode()
		context.Logger.Debugf("Register: token = %s, verificationCode = %s", token, verificationCode)

		context.Storage.RegistrationDAO().Upsert(
			token,
			&models.User{
				Login:     request.Login,
				Password:  request.Password,
				Email:     request.Email,
				FirstName: request.FirstName,
				LastName:  request.LastName,
			},
			time.Now().Add(internal_data.TokenTTL),
			verificationCode,
		)

		go func() {
			if err := context.Mailer.SendEmail([]string{request.Email}, notifier.BuildMessage(verificationCode,
				request.FirstName, request.LastName)); err != nil {
				context.Logger.Errorf("Register: can't send email to %s, %s", request.Email, err)
			}
		}()

		return http.StatusOK, &RegisterResponse{
			Response: &RegisterWrapper{
				RandomToken: token,
				ExpiresIn:   strconv.Itoa(int(internal_data.ResentCodeIn.Seconds())),
			},
		}
	}

	return handleAccountsRegister(request)
}
