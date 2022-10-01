package accounts

import (
	"encoding/hex"
	"encoding/json"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage"
	"github.com/sava-cska/SPbSU-EMKN/internal/utils"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func HandleAccountsRegister(logger *logrus.Logger, storage *storage.Storage) http.HandlerFunc {
	expireIn := 60 * time.Second
	verificationCodeLength := 6
	tokenLength := 20

	const (
		loginIsNotAvailable int = iota
		illegalPassword
		illegalLogin
		illegalEmail
	)

	validate := func(request *RegisterRequest) (int, *RegisterErrors) {
		if len(request.Login) == 0 {
			return http.StatusBadRequest, &RegisterErrors{
				IllegalLogin: &Error{Code: illegalLogin},
			}
		}
		if len(request.Password) == 0 {
			return http.StatusBadRequest, &RegisterErrors{
				IllegalPassword: &Error{Code: illegalPassword},
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
		if storage.UserDao().Exists(request.Login) {
			return http.StatusBadRequest, &RegisterResponse{
				Errors: &RegisterErrors{LoginIsNotAvailable: &Error{Code: loginIsNotAvailable}},
			}
		}
		token := generateToken()
		verificationCode := generateVerificationCode()
		_ = storage.RegistrationDao().Upsert(
			token,
			request.Login,
			request.Password,
			request.Email,
			request.FirstName,
			request.LastName,
			time.Now().Add(expireIn),
			verificationCode,
		)

		if err := utils.SendEmail(request.Email, verificationCode, request.FirstName, request.LastName); err != nil {
			return http.StatusBadRequest, &RegisterResponse{
				Errors: &RegisterErrors{IllegalEmail: &Error{Code: illegalEmail}},
			}
		}

		return http.StatusOK, &RegisterResponse{
			Response: &WrapResponse{
				RandomToken: token,
				ExpiresIn:   strconv.Itoa(int(expireIn.Seconds())),
			},
		}
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		logger.Debugf("HandleAccountsRegister - Called URI %s", request.RequestURI)

		var registerRequest RegisterRequest
		if errJSON := utils.ParseBody(interface{}(&registerRequest), request); errJSON != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write([]byte("Can't parse query"))
			return
		}

		code, resp := handleAccountsRegister(&registerRequest)
		writer.WriteHeader(code)

		respJSON, errRespJSON := json.Marshal(resp)
		if errRespJSON != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write([]byte("Can't create JSON object from data."))
			return
		}

		_, _ = writer.Write(respJSON)
	}
}
