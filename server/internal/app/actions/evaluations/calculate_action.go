package evaluations

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Knetic/govaluate"
	"github.com/sava-cska/SPbSU-Calculator/internal/app/storage"
	"github.com/sava-cska/SPbSU-Calculator/internal/utils"
	"github.com/sirupsen/logrus"
)

func HandleEvaluationsCalculate(logger *logrus.Logger, storage *storage.Storage) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		logger.Debugf("HandleEvaluationsCalculate - Called URI %s", request.RequestURI)

		var expression CalculateRequest
		if errJSON := utils.ParseBody(logger, interface{}(&expression), writer, request); errJSON != nil {
			utils.HandleError(logger, writer, http.StatusBadRequest, "Can't parse Calculate query", errJSON)
			return
		}

		defer func() {
			if r := recover(); r != nil {
				utils.HandleError(logger, writer, http.StatusBadRequest, r.(string), nil)
				storage.Evaluations().Upsert(expression.UserUid,
					encode(expression.Evaluation, logger, writer),
					nil) // use nil if you need to upsert evaluation with error
			}
		}()

		stringExpr := toString(expression.Evaluation)
		formula, errCreateFormula := govaluate.NewEvaluableExpression(stringExpr)
		if errCreateFormula != nil {
			panic(fmt.Sprintf("Can't create formula from string representation %s. %s", stringExpr, errCreateFormula))
		}

		rawResult, errEvaluate := formula.Evaluate(nil)
		if errEvaluate != nil {
			panic(fmt.Sprintf("Can't evaluate formula %s. %s", stringExpr, errEvaluate))
		}

		result := fmt.Sprintf("%v", rawResult)
		logger.Debugf("Result of the expression %s equals %s", stringExpr, result)

		storage.Evaluations().Upsert(expression.UserUid,
			encode(expression.Evaluation, logger, writer),
			&result, // use nil if you need to upsert evaluation with error
		)

		postResponse := CalculateResponse{Result: result}
		respJSON, errRespJSON := json.Marshal(postResponse)
		if errRespJSON != nil {
			utils.HandleError(logger, writer, http.StatusInternalServerError, "Can't create JSON object from data.", errRespJSON)
			return
		}

		writer.Write(respJSON)
	}
}
