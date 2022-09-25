package evaluations

import (
	"encoding/json"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/storage"
	"github.com/sava-cska/SPbSU-EMKN/internal/utils"
	"net/http"

	"github.com/sirupsen/logrus"
)

func HandleEvaluationsList(logger *logrus.Logger, storage *storage.Storage) http.HandlerFunc {

	return func(writer http.ResponseWriter, request *http.Request) {
		logger.Debugf("HandleEvaluationsList - Called URI %s", request.RequestURI)
		var user ListRequest
		if errJSON := utils.ParseBody(logger, interface{}(&user), writer, request); errJSON != nil {
			utils.HandleError(logger, writer, http.StatusBadRequest, "Can't parse List query", errJSON)
			return
		}

		evals, errList := storage.Evaluations().List(user.UserUid)
		if errList != nil {
			utils.HandleError(logger, writer, http.StatusInternalServerError, "Can't get list of evaluations.", errList)
			return
		}

		evaluationHistory := make([]EvaluationHistory, len(evals))
		for idx := 0; idx < len(evals); idx++ {
			evaluationHistory[idx] = EvaluationHistory{
				Evaluation: decode(evals[idx].Evaluation, logger, writer),
				Result:     evals[idx].Result,
			}
		}

		response := ListResponse{
			EvaluationHistory: evaluationHistory,
		}

		respJSON, errRespJSON := json.Marshal(response)
		if errRespJSON != nil {
			utils.HandleError(logger, writer, http.StatusInternalServerError, "Can't create JSON object from data.", errRespJSON)
			return
		}

		writer.Write(respJSON)
	}
}
