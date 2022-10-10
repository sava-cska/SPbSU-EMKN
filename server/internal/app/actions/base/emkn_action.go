package base

import (
	"encoding/json"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/core/dependency"
	"github.com/sava-cska/SPbSU-EMKN/internal/utils"
	"net/http"
)

func HandleAction[Req Request, Res Response](
	path string,
	businessLogicHandler func(*Req, *dependency.DependencyContext) (int, *Res),
	context *dependency.DependencyContext,
) {
	handlerFunc := func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		context.Logger.Debugf("HandleAction - Called URI %s", httpRequest.RequestURI)

		var registerRequest Req
		if errJSON := utils.ParseBody(interface{}(&registerRequest), httpRequest); errJSON != nil {
			utils.HandleError(context.Logger, responseWriter, http.StatusBadRequest, "Can't parse httpRequest.", errJSON)
			return
		}

		code, resp := businessLogicHandler(&registerRequest, context)

		respJSON, errRespJSON := json.Marshal(resp)
		if errRespJSON != nil {
			utils.HandleError(context.Logger, responseWriter, http.StatusInternalServerError, "Can't create JSON object from data.", errRespJSON)
			return
		}

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(code)
		_, _ = responseWriter.Write(respJSON)
	}
	context.Router.HandleFunc(path, handlerFunc)
}
