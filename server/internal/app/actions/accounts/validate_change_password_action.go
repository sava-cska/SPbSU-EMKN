package accounts

import (
	"encoding/hex"
	"encoding/json"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage"
	"github.com/sava-cska/SPbSU-EMKN/internal/utils"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"time"
)

func HandleValidateChangePassword(logger *logrus.Logger, storage *storage.Storage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		logger.Debugf("HandleAccountsValidateChangePassword - Called URI %s", request.RequestURI)

		var parsedRequest ValidateChangePasswordRequest
		if err := utils.ParseBody(interface{}(&parsedRequest), request); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte("Failed to parse request body"))
			return
		}

		correctVerificationCode, expiresAt, err := storage.ChangePasswordDao().GetVerificationCodeInfo(parsedRequest.RandomToken)
		if err != nil {
			utils.HandleError(logger, writer, http.StatusInternalServerError, "Failed to get verification code", err)
			return
		}
		if correctVerificationCode == "" {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte("Failed to find random token"))
			return
		}

		errors := validateVerificationCode(correctVerificationCode, parsedRequest.VerificationCode, expiresAt)

		var responseBody ValidateChangePasswordResponse
		if errors == nil {
			token := generateToken(20)
			err = storage.ChangePasswordDao().SetChangePasswordToken(parsedRequest.RandomToken, token)
			if err != nil {
				utils.HandleError(logger, writer, http.StatusInternalServerError, "Failed to store changePasswordToken", err)
				return
			}
			responseBody = ValidateChangePasswordResponse{
				ChangePasswordToken: token,
			}
			writer.WriteHeader(http.StatusOK)
		} else {
			responseBody = ValidateChangePasswordResponse{
				Errors: errors,
			}
			writer.WriteHeader(http.StatusBadRequest)
		}

		body, err := json.Marshal(&responseBody)
		_, _ = writer.Write(body)
	}
}

func generateToken(length uint16) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func validateVerificationCode(correctVerificationCode string, verificationCode string, expiresAt *time.Time) *ErrorsUnion{
	if expiresAt.Before(time.Now()) {
		return &ErrorsUnion{ ChangePasswordExpired: &Error{} }
	}

	if correctVerificationCode != verificationCode {
		return &ErrorsUnion{ InvalidCode: &Error{} }
	}

	return nil
}

type ValidateChangePasswordRequest struct {
	RandomToken string `json:"random_token"`
	VerificationCode string `json:"verification_code"`
}

type ValidateChangePasswordResponse struct {
	ChangePasswordToken string `json:"change_password_token,omitempty"`
	Errors *ErrorsUnion `json:"errors,omitempty"`
}
