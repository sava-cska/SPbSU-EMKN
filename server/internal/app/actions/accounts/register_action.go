package accounts

import (
	"encoding/json"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/actions"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage"
	"github.com/sava-cska/SPbSU-EMKN/internal/utils"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"strconv"
)

func HandleAccountsRegister(logger *logrus.Logger, storage *storage.Storage) http.HandlerFunc {
	handleAccountsRegister := func(request RegisterRequest) RegisterResponse {
		errors := make([]actions.Error, 0)
		r := rand.Float32()
		switch {
		case 0 <= r && r <= 0.2:
			errors = append(errors, actions.Error{Code: actions.IllegalEmail})
		case 0.2 < r && r <= 0.4:
			errors = append(errors, actions.Error{Code: actions.LoginIsNotAvailable})
		case 0.4 < r && r <= 0.6:
			errors = append(errors, actions.Error{Code: actions.IllegalPassword})
		}
		return RegisterResponse{
			EmknResponse: actions.EmknResponse{
				Errors: errors,
			},
			RandomToken: strconv.Itoa(rand.Int()),
			ExpiresIn:   "150",
		}
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		logger.Debugf("HandleAccountsRegister - Called URI %s", request.RequestURI)

		var registerRequest RegisterRequest
		if errJSON := utils.ParseBody(logger, interface{}(&registerRequest), writer, request); errJSON != nil {
			utils.HandleError(logger, writer, http.StatusBadRequest, "Can't parse Calculate query", errJSON)
			return
		}

		respJSON, errRespJSON := json.Marshal(handleAccountsRegister(registerRequest))
		if errRespJSON != nil {
			utils.HandleError(logger, writer, http.StatusInternalServerError, "Can't create JSON object from data.", errRespJSON)
			return
		}

		writer.Write(respJSON)
	}
}
