package accounts

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage"
	"github.com/sava-cska/SPbSU-EMKN/internal/utils"
	"github.com/sirupsen/logrus"
)

func checkCorrectnessCode(expireTime time.Time, currentTime time.Time, verificationCodeFromDB string,
	verificationCodeFromRequest string) (ValidateEmailResponse, error) {
	if expireTime.Before(currentTime) {
		return ValidateEmailResponse{Errors: &ErrorsUnion{RegistrationExpired: &Error{}}}, errors.New("time expired")
	}

	if verificationCodeFromDB != verificationCodeFromRequest {
		return ValidateEmailResponse{Errors: &ErrorsUnion{InvalidCode: &Error{}}},
			errors.New("verification code isn't correct")
	}
	return ValidateEmailResponse{}, nil
}

func HandleAccountsValidateEmail(logger *logrus.Logger, storage *storage.Storage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		logger.Debugf("HandleAccountsValidateEmail - Called URI %s", request.RequestURI)

		currentTime := time.Now()
		var validationRequest ValidateEmailRequest
		if errJSON := utils.ParseBody(interface{}(&validationRequest), request); errJSON != nil {
			utils.HandleError(logger, writer, http.StatusBadRequest, "Can't parse /accounts/validate_email request.", errJSON)
			return
		}

		logger.Debugf("ValidateEmail: token = %s, code = %s", validationRequest.Token, validationRequest.VerificationCode)

		user, expireTime, verificationCode, err := storage.RegistrationDAO().FindRegistration(validationRequest.Token)
		if err != nil {
			utils.HandleError(logger, writer, http.StatusBadRequest, "Can't find registration record.", err)
			return
		}

		response, err := checkCorrectnessCode(expireTime, currentTime, verificationCode, validationRequest.VerificationCode)
		responseJSON, errJSON := json.Marshal(response)
		if errJSON != nil {
			utils.HandleError(logger, writer, http.StatusInternalServerError, "Can't create JSON object from data.", errJSON)
			return
		}

		var returnCode int
		if err == nil {
			returnCode = http.StatusOK
			if errAdd := storage.UserDAO().AddUser(&user); errAdd != nil {
				utils.HandleError(logger, writer, http.StatusInternalServerError, "Can't add user into database.", errAdd)
				return
			}
		} else {
			returnCode = http.StatusBadRequest
		}

		writer.WriteHeader(returnCode)
		writer.Write([]byte(responseJSON))
	}
}
