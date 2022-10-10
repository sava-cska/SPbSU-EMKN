package base

import (
	"encoding/json"
	"github.com/gdexlab/go-render/render"
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
		context.Logger.Debugf("%s - Called URI %s", path, httpRequest.RequestURI)

		var request Req
		if errJSON := utils.ParseBody(interface{}(&request), httpRequest); errJSON != nil {
			utils.HandleError(context.Logger, responseWriter, http.StatusBadRequest, "Can't parse httpRequest.", errJSON)
			return
		}

		context.Logger.Debugf("%s\n\trequest: %s", path, render.AsCode(request))

		code, response := businessLogicHandler(&request, context)

		respJSON, errRespJSON := json.Marshal(response)
		if errRespJSON != nil {
			utils.HandleError(context.Logger, responseWriter, http.StatusInternalServerError, "Can't create JSON object from data.", errRespJSON)
			return
		}
		context.Logger.Debugf("%s\n\tcode: %d\n\tresponse: %s", path, code, render.AsCode(response))

		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(code)
		_, _ = responseWriter.Write(respJSON)
	}
	context.Router.HandleFunc(path, handlerFunc)
}
