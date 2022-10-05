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

func checkExpiredTimeAndNewPwd(changePwdExpiredTime time.Time, currentTime time.Time,
	newPassword string) (CommitPwdChangeResponse, error) {
	if changePwdExpiredTime.Before(currentTime) {
		response := CommitPwdChangeResponse{Errors: &ErrorsUnion{ChangePasswordExpired: &Error{}}}
		return response, errors.New("time is expired")
	}

	if len(newPassword) == 0 {
		response := CommitPwdChangeResponse{Errors: &ErrorsUnion{IllegalPassword: &Error{}}}
		return response, errors.New("empty password")
	}

	return CommitPwdChangeResponse{}, nil
}

func commitPwdChange(logger *logrus.Logger, writer http.ResponseWriter, storage *storage.Storage,
	requestCommit *CommitPwdChangeRequest, currentTime time.Time) (CommitPwdChangeResponse, int, error) {
	login, changePwdExpiredTime, errDB := storage.ChangePasswordDao().FindPwdToken(requestCommit.ChangePwdToken)
	if errDB != nil {
		return CommitPwdChangeResponse{}, 0, errDB
	}

	response, errCheck := checkExpiredTimeAndNewPwd(changePwdExpiredTime, currentTime, requestCommit.NewPassword)
	if errCheck != nil {
		return response, http.StatusInternalServerError, nil
	}

	errUpdate := storage.UserDAO().UpdatePassword(login, requestCommit.NewPassword)
	if errUpdate != nil {
		return CommitPwdChangeResponse{}, 0, errUpdate
	}

	return CommitPwdChangeResponse{}, http.StatusOK, nil
}

func HandleAccountsCommitPwdChange(logger *logrus.Logger, storage *storage.Storage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		logger.Debugf("HandleAccountsCommitPwdChange - Called URI %s", request.RequestURI)

		currentTime := time.Now()

		var commitPwdChangeRequest CommitPwdChangeRequest
		if errJSON := utils.ParseBody(interface{}(&commitPwdChangeRequest), request); errJSON != nil {
			utils.HandleError(logger, writer, http.StatusBadRequest, "Can't parse /accounts/commit_change_password request.", errJSON)
			return
		}

		response, returnCode, err := commitPwdChange(logger, writer, storage, &commitPwdChangeRequest, currentTime)
		if err != nil {
			utils.HandleError(logger, writer, http.StatusInternalServerError, "Can't update password.", err)
			return
		}

		responseJSON, errJSON := json.Marshal(response)
		if errJSON != nil {
			utils.HandleError(logger, writer, http.StatusInternalServerError, "Can't create JSON object from data.", errJSON)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(returnCode)
		writer.Write(responseJSON)
	}
}
