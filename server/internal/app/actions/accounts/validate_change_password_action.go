package accounts

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage"
	"github.com/sava-cska/SPbSU-EMKN/internal/utils"
	"github.com/sirupsen/logrus"
)

func HandleValidateChangePassword(logger *logrus.Logger, storage *storage.Storage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		logger.Debugf("HandleAccountsValidateChangePassword - Called URI %s", request.RequestURI)

		var parsedRequest = ValidateChangePasswordRequest{}
		if err := utils.ParseBody(interface{}(&parsedRequest), request); err != nil {
			utils.HandleError(logger, writer, http.StatusBadRequest, "Failed to parse request body", err)
			return
		}

		correctVerificationCode, expiresAt, err := storage.ChangePasswordDao().GetVerificationCodeInfo(parsedRequest.RandomToken)
		if err != nil {
			utils.HandleError(logger, writer, http.StatusInternalServerError, "Failed to get verification code", err)
			return
		}
		if correctVerificationCode == "" {
			utils.HandleError(logger, writer, http.StatusBadRequest, "Failed to find random token", nil)
			return
		}

		errors := validateVerificationCode(correctVerificationCode, parsedRequest.VerificationCode, expiresAt)

		var responseBody ValidateChangePasswordResponse
		var statusCode int
		if errors == nil {
			token := utils.GenerateToken()
			err = storage.ChangePasswordDao().SetChangePasswordToken(parsedRequest.RandomToken, token)
			if err != nil {
				utils.HandleError(logger, writer, http.StatusInternalServerError, "Failed to store changePasswordToken", err)
				return
			}
			responseBody = ValidateChangePasswordResponse{
				ChangePasswordToken: token,
			}
			statusCode = http.StatusOK
		} else {
			responseBody = ValidateChangePasswordResponse{
				Errors: errors,
			}
			statusCode = http.StatusBadRequest
		}

		body, err := json.Marshal(responseBody)
		if err != nil {
			utils.HandleError(logger, writer, http.StatusInternalServerError, "Can't create JSON object from data.", err)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(statusCode)
		writer.Write(body)
	}
}

func validateVerificationCode(correctVerificationCode string, verificationCode string, expiresAt *time.Time) *ErrorsUnion {
	if expiresAt.Before(time.Now()) {
		return &ErrorsUnion{ChangePasswordExpired: &Error{}}
	}

	if correctVerificationCode != verificationCode {
		return &ErrorsUnion{InvalidCode: &Error{}}
	}

	return nil
}
