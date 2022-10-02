package accounts

import (
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage"
	"github.com/sava-cska/SPbSU-EMKN/internal/utils"
	"github.com/sirupsen/logrus"
)

func HandleAccountsRegister(logger *logrus.Logger, storage *storage.Storage) http.HandlerFunc {
	expireIn := 60 * time.Second
	verificationCodeLength := 6
	tokenLength := 20

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
		return http.StatusOK, nil
	}

	generateToken := func() string {
		b := make([]byte, tokenLength)
		if _, err := rand.Read(b); err != nil {
			return ""
		}
		return hex.EncodeToString(b)
	}

	generateVerificationCode := func() string {
		code := strings.Builder{}
		for i := 0; i < verificationCodeLength; i++ {
			code.WriteString(strconv.Itoa(rand.Intn(10)))
		}
		return code.String()
	}

	handleAccountsRegister := func(request *RegisterRequest) (int, *RegisterResponse) {
		if code, errors := validate(request); errors != nil {
			return code, &RegisterResponse{
				Errors: errors,
			}
		}
		if storage.UserDAO().Exists(request.Login) {
			return http.StatusBadRequest, &RegisterResponse{
				Errors: &ErrorsUnion{LoginIsNotAvailable: &Error{}},
			}
		}
		token := generateToken()
		verificationCode := generateVerificationCode()
		storage.RegistrationDAO().Upsert(
			token,
			request,
			time.Now().Add(expireIn),
			verificationCode,
		)

		if err := utils.SendEmail(request.Email, verificationCode, request.FirstName, request.LastName); err != nil {
			return http.StatusBadRequest, &RegisterResponse{
				Errors: &ErrorsUnion{IllegalEmail: &Error{}},
			}
		}

		return http.StatusOK, &RegisterResponse{
			Response: &RegisterWrapper{
				RandomToken: token,
				ExpiresIn:   strconv.Itoa(int(expireIn.Seconds())),
			},
		}
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		logger.Debugf("HandleAccountsRegister - Called URI %s", request.RequestURI)

		var registerRequest RegisterRequest
		if errJSON := utils.ParseBody(interface{}(&registerRequest), request); errJSON != nil {
			utils.HandleError(logger, writer, http.StatusBadRequest, "Can't parse /accounts/register request.", errJSON)
			return
		}

		code, resp := handleAccountsRegister(&registerRequest)

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
